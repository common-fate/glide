import {
  ArrowBackIcon,
  CheckCircleIcon,
  CheckIcon,
  ChevronDownIcon,
  ChevronRightIcon,
  WarningIcon,
} from "@chakra-ui/icons";
import {
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Badge,
  Box,
  Button,
  Center,
  Circle,
  CircularProgress,
  Code,
  Container,
  Flex,
  FormControl,
  FormHelperText,
  FormLabel,
  Grid,
  GridItem,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverHeader,
  PopoverTrigger,
  useClipboard,
  Spinner,
  Stack,
  Text,
  UnorderedList,
  OrderedList,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import Confetti from "react-confetti";
import { Helmet } from "react-helmet";
import { Controller, useForm } from "react-hook-form";
import { Link, useMatch, useNavigate } from "react-location";
import ReactMarkdown from "react-markdown";
import { Sticky, StickyContainer } from "react-sticky";
import useWindowSize from "react-use/lib/useWindowSize";
import {
  CFCode,
  CFReactMarkownCode,
} from "../../../components/CodeInstruction";
import { ConnectorArrow } from "../../../components/ConnectorArrow";
import { ExpandingImage } from "../../../components/ExpandingImage";
import { CommonFateLogo } from "../../../components/icons/Logos";
import { ProviderIcon } from "../../../components/icons/providerIcon";
import { UserLayout } from "../../../components/Layout";
import { useGetProvider } from "../../../utils/backend-client/registry/orval";

const Page = () => {
  const {
    params: { id },
  } = useMatch();

  const idDecoded = decodeURIComponent(id);
  // break the the id into `name` and `team`
  const [name, team] = idDecoded.split("/");

  // now we can use the name and team to look up the provider
  console.log("👀 name, team", name, team);

  // TODO:
  // fetch data and hydrate UI with relevant details... we'll add more detail at a later stage
  // const { data } = useGetRegistryProvider({ name, team, version: "latest" });

  const navigate = useNavigate();
  const { width, height } = useWindowSize();

  // const [validationErrorMsg, setValidationErrorMsg] = useState("");
  // const { hasCopied, onCopy } = useClipboard(validationErrorMsg);

  return (
    <UserLayout>
      <Stack
        justifyContent={"center"}
        alignItems={"center"}
        spacing={{ base: 1, md: 0 }}
        borderBottom="1px solid"
        borderColor="neutrals.200"
        h="80px"
        py={{ base: 4, md: 0 }}
        flexDirection={{ base: "column", md: "row" }}
      >
        <IconButton
          as={Link}
          aria-label="Go back"
          pos="absolute"
          left={4}
          icon={<ArrowBackIcon />}
          rounded="full"
          variant="ghost"
          to="/"
        />
      </Stack>
    </UserLayout>
  );
};

export default Page;
