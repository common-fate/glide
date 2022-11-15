import { ArrowBackIcon, LockIcon } from "@chakra-ui/icons";
import {
  Center,
  Container,
  HStack,
  IconButton,
  Stack,
  Text,
} from "@chakra-ui/react";
import { Link, MakeGenerics, useMatch, useSearch } from "react-location";
import { UserLayout } from "../../../components/Layout";

import { useCallback, useEffect, useRef, useState } from "react";
import { Helmet } from "react-helmet";
import useWebSocket from "react-use-websocket";
import { Terminal } from "xterm";
import { callRequestOperation } from "../../../utils/backend-client/end-user/end-user";
import { useUser } from "../../../utils/context/userContext";
import "./xterm.css";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    peer?: string;
  };
}>;

interface TtyOutputMessage {
  t: "tty_output";
  c: {
    data: string;
  };
}

interface SessionUpdatedMessage {
  t: "session_updated";
  c: {
    status: string;
  };
}

type Message = TtyOutputMessage | SessionUpdatedMessage;

const Home = () => {
  const {
    params: { id: requestId },
  } = useMatch();
  const user = useUser();
  const search = useSearch<MyLocationGenerics>();
  useEffect(() => {}, []);

  const getSocketUrl = useCallback(() => {
    return new Promise<string>(async (resolve) => {
      const res = await callRequestOperation(requestId, {
        operation: "get-socket",
      });
      const url = res.data.url as string;
      resolve(url);
    });
  }, [requestId]);

  const { sendMessage, getWebSocket } = useWebSocket(getSocketUrl, {
    onOpen: () => console.log("opened"),
    //Will attempt to reconnect on all close events, such as server shutting down
    shouldReconnect: (closeEvent) => true,
    onMessage: (event) => {
      const msg = JSON.parse(event.data) as Message;
      if (msg.t === "session_updated") {
        xtermRef.current?.writeln(`session is now ${msg.c.status}`);
      }
      if (msg.t === "tty_output") {
        xtermRef.current?.writeln(msg.c.data);
        xtermRef.current?.write("> ");
      }
    },
  });

  const ws = getWebSocket();

  const xtermRef = useRef<Terminal>();
  const [input, setInput] = useState("");

  useEffect(() => {
    if (xtermRef.current == null && ws != null) {
      const el = document.getElementById("xterm");
      if (el != null) {
        xtermRef.current = new Terminal({
          convertEol: true,
          rows: 40,
          cols: 120,
          theme: {
            background: "#2D2F30",
          },
        });
        xtermRef.current.open(el);
        xtermRef.current.writeln("opening a session...");
      } else {
        console.error("no xterm id found");
      }
    }
  }, [xtermRef.current, ws]);

  const handleMessage = (inputToSend: string) => {
    if (search.peer != null) {
      xtermRef.current?.writeln("");
      xtermRef.current?.writeln(
        `command is pending: '${inputToSend}' - awaiting a peer to join your session to approve command`
      );
      xtermRef.current?.write("> ");
    } else {
      const command = {
        t: "command",
        c: {
          command: inputToSend,
        },
      };
      const commandString = JSON.stringify(command);
      sendMessage(commandString);
      xtermRef.current?.writeln("");
    }
  };

  useEffect(() => {
    if (xtermRef.current) {
      const token = xtermRef.current.onData((data: string) => {
        const code = data.charCodeAt(0);
        // If the user hits empty and there is something typed echo it.
        if (code === 13 && input.length > 0) {
          handleMessage(input);
          setInput("");
        } else if (code < 32) {
          // Disable control Keys such as arrow keys
          return;
        } else if (code === 127) {
          //backspace
          xtermRef.current?.write("\b \b");
          setInput(input.slice(0, -1));
        } else {
          // Add general key press characters to the terminal
          xtermRef.current?.write(data);
          setInput(input + data);
        }
      });
      return () => token.dispose();
    }
  }, [input, xtermRef.current, xtermRef, ws]);

  return (
    <div>
      <UserLayout>
        <Helmet>
          <title>Web Shell</title>
        </Helmet>
        {/* The header bar */}
        <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
          <IconButton
            as={Link}
            aria-label="Go back"
            pos="absolute"
            left={4}
            icon={<ArrowBackIcon />}
            rounded="full"
            variant="ghost"
            to={"/requests/" + requestId}
          />
          <HStack>
            <Text as="h4" textStyle="Heading/H4">
              Web Shell {search.peer != null && "- Peer Review Mode"}
            </Text>
            {search.peer != null && <LockIcon />}
          </HStack>
        </Center>
        {/* Main content */}
        <Container maxW="container.xl" py={16}>
          <Stack
            bg="neutrals.700"
            w="1150px"
            // h="700px"
            py={8}
            borderWidth={"1px"}
            justifyContent="center"
            alignItems={"center"}
            borderRadius="8px"
            id="xterm"
          ></Stack>
        </Container>
      </UserLayout>
    </div>
  );
};

export default Home;
