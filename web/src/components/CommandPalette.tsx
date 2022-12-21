import React, { useEffect, useState } from "react";
import { matchSorter } from "match-sorter";
import {
  Badge,
  Box,
  Button,
  chakra,
  Flex,
  Highlight,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Kbd,
  Modal,
  ModalBody,
  ModalContent,
  ModalHeader,
  Text,
  ModalOverlay,
  ModalProps,
  Spinner,
  Tag,
  Tooltip,
  IconProps,
  ComponentWithAs,
  useUpdateEffect,
} from "@chakra-ui/react";
import { useRef } from "react";
import {
  AddIcon,
  ArrowForwardIcon,
  ArrowRightIcon,
  CheckIcon,
  DeleteIcon,
  EditIcon,
  PlusSquareIcon,
  StarIcon,
} from "@chakra-ui/icons";
import { useNavigate, useRouter } from "react-location";
import { GitCompareOutline } from "./icons/Icons";
import { useListUserAccessRules } from "../utils/backend-client/end-user/end-user";
import { ProviderIcon } from "./icons/providerIcon";
import {
  userListFavorites,
  useUserListFavorites,
} from "../utils/backend-client/default/default";

interface Props {}

type ICommand = {
  name: string;
  icon: ComponentWithAs<"svg", IconProps>;
  action: () => void;
  isAdminOnly?: boolean;
  type?: "favorite";
};

const EditActionModal = ({
  isOpen,
  onClose,
  ...rest
}: Omit<ModalProps, "children">) => {
  const [selectedKey, setSelectedKey] = useState("s3");

  const [inputValue, setInputValue] = useState("");

  const inputRef = useRef(null);

  const nav = useNavigate();

  const otherActions: ICommand[] = [
    {
      name: "Create access rule",
      icon: ArrowForwardIcon,
      action: () => nav({ to: "/admin/access-rules/create" }),
      isAdminOnly: true,
    },
  ];

  /**
   Context aware
    - The palette knows where you are in the app
    - i.e. If you open the palette on an Access Request it displays a set of actions specific to that request e.g. 'Approve' or 'Revoke' request, or 'Copy access instructions'
   */
  // const accessRequestActions: ICommand[] = [
  //   {
  //     name: "Approve request",
  //     icon: CheckIcon,
  //     action: () => undefined,
  //   },
  //   {
  //     name: "Revoke request",
  //     icon: DeleteIcon,
  //     action: () => undefined,
  //   },
  // ];

  type ContextualCommand = {
    pathRegex: string;
    commands: ICommand[];
  };

  const { data: rules } = useListUserAccessRules();

  const rulesAsCommands: ICommand[] = rules?.accessRules
    ? rules?.accessRules.map((rule) => ({
        name: rule.name,
        action: () => nav({ to: `/access/request/${rule.id}` }),
        icon: (props) => (
          <ProviderIcon
            shortType={rule.target.provider.type}
            h="4"
            w="4"
            {...props}
          />
        ),
      }))
    : [];

  // now do the same for favourites
  const { data: favorites } = useUserListFavorites();

  const favoritesAsCommands: ICommand[] = favorites?.favorites
    ? favorites?.favorites.map((favorite) => ({
        name: favorite.name,
        action: () =>
          nav({
            to: `/access/request/${favorite.ruleId}?favorite=${favorite.id}`,
          }),
        icon: StarIcon,
        type: "favorite",
      }))
    : [];

  const contextualCommands: ContextualCommand[] = [
    {
      pathRegex: "/admin",
      commands: [
        {
          name: "Switch to user",
          icon: GitCompareOutline,
          action: () => nav({ to: "/" }),
          isAdminOnly: true,
        },
      ],
    },
    {
      // regex pattern for paths that dont include '/admin'
      pathRegex: "(?!/admin)",
      commands: [
        {
          name: "Switch to admin",
          icon: GitCompareOutline,
          action: () => nav({ to: "/admin/" }),
          isAdminOnly: true,
        },
      ],
    },
  ];

  /**
   * To handle contextual command palette results we need to update the input arrays (ICommand) depending on the page route,
   * we need to watch this route and ensure newly updated values are passed to the results filtering function
   *
   * Another way we could handle contextual results (including support for arbitrary/custom conditions) is:
   * - Define a set of 'contexts' which are a set of conditions that must be met for a set of commands to be displayed
   * - Define a set of 'commands' which are a set of commands that can be displayed in the palette
   * - Define a set of 'rules' which are a set of conditions that must be met for a set of commands to be displayed
   */

  const router = useRouter();

  const contextualAndStaticCommands = React.useMemo(() => {
    // console.log({ router });
    // how do I get the router current path using react-location and not using window object
    const currentPath = router.state.location.href || "";
    const matchedContextualCommands = contextualCommands
      // run the regex from string to regex
      .filter((command) => new RegExp(command.pathRegex).test(currentPath))
      .map((command) => command.commands)
      .flat();
    return [
      ...rulesAsCommands,
      ...favoritesAsCommands,
      // ...matchedContextualCommands,
      // ...otherActions,
    ];
  }, [router.state.location.href, rulesAsCommands, favoritesAsCommands]);

  const results = React.useMemo(
    function getResults() {
      if (inputValue.length < 2) return contextualAndStaticCommands;
      return matchSorter(contextualAndStaticCommands, inputValue, {
        keys: ["name"],
      }).slice(0, 20);
    },
    [inputValue, contextualAndStaticCommands, router.state.location.href]
  );

  const [active, setActive] = React.useState(0);
  const menuRef = React.useRef<HTMLDivElement>(null);
  const eventRef = React.useRef<"mouse" | "keyboard">("keyboard");

  const onKeyUp = React.useCallback((e: React.KeyboardEvent) => {
    eventRef.current = "keyboard";
    switch (e.key) {
      case "Control":
      case "Alt":
      case "Shift": {
        e.preventDefault();
      }
    }
  }, []);

  const onKeyDown = React.useCallback(
    (e: React.KeyboardEvent) => {
      eventRef.current = "keyboard";
      switch (e.key) {
        case "ArrowDown": {
          e.preventDefault();
          if (active + 1 < results.length) {
            setActive(active + 1);
          }
          break;
        }
        case "ArrowUp": {
          e.preventDefault();
          if (active - 1 >= 0) {
            setActive(active - 1);
          }
          break;
        }
        case "Control":
        case "Alt":
        case "Shift": {
          e.preventDefault();
          onClose();
          break;
        }
        case "Enter": {
          if (results?.length <= 0) {
            break;
          }
          onClose();
          results[active].action();
          // also clear the input value and reset active
          setInputValue("");
          setActive(0);
          break;
        }
      }
    },
    [active, results, nav]
  );

  useUpdateEffect(() => {
    setActive(0);
  }, [inputValue]);

  /** used to scroll input results into view */
  // useUpdateEffect(() => {
  //   if (!menuRef.current || eventRef.current === "mouse") return;
  //   const node = menuNodes.map.get(active);
  //   if (!node) return;
  // scrollIntoView(node, {
  //   scrollMode: "if-needed",
  //   block: "nearest",
  //   inline: "nearest",
  //   boundary: menuRef.current,
  // });
  // }, [active]);

  const loading = false;

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      size="lg"
      onCloseComplete={() => {
        setInputValue("");
      }}
    >
      {/* <ModalOverlay /> */}
      <ModalContent
        border="1px solid"
        borderColor="whiteAlpha.400"
        overflow="hidden"
        bg="#ffffff76"
        rounded="md"
        backdropFilter="blur(20px) saturate(170%) contrast(50%) brightness(130%)"
        boxShadow="rgb(0 0 0 / 50%) 0px 16px 70px"
      >
        {/* <ModalCloseButton zIndex={999} size="sm" /> */}
        {/* <ModalHeader fontSize="md" pb={2}>
          Add an action
        </ModalHeader> */}
        <ModalBody
          p={0}
          position="relative"
          pb={3}
          h="100%"
          maxH="80vh"
          ref={menuRef}
        >
          <Flex flex={1} position="relative" flexDir="column">
            <InputGroup pt={4}>
              <Input
                pb={5}
                // To increase the placeholder text size you can use the `fontSize` prop
                fontSize="xl"
                // py={5}
                // minH={20}
                spellCheck={false}
                px={6}
                variant="flushed"
                size="lg"
                onChange={(e) => setInputValue(e.target.value)}
                value={inputValue}
                // onKeyPress={(e) => {
                //   e.key === ":" && keyCheck();
                //   e.key === "Enter" && keyCheck();
                // }}
                autoFocus={true}
                type="text"
                ref={inputRef}
                placeholder="Type a command or search"
                onKeyDown={onKeyDown}
                onKeyUp={onKeyUp}
              />
              {loading && (
                <InputRightElement>
                  <Spinner size="sm" />
                </InputRightElement>
              )}
            </InputGroup>
            {/* Iteration of all privileges (search results) */}
            <Box
              flex="1 0 auto"
              overflowY="auto"
              maxH="60vh"
              py={3}
              sx={{
                "&::-webkit-scrollbar": {
                  WebkitAppearance: "none",
                  width: "7px",
                },
                "&::-webkit-scrollbar-thumb": {
                  borderRadius: "4px",
                  backgroundColor: "rgba(0, 0, 0, .3)",
                  boxShadow: "0 0 1px rgba(255, 255, 255, .3)",
                },
              }}
            >
              {results.map((el, index) => {
                const selected = index === active;
                // const isLvl1 = item.type === 'lvl1'
                return (
                  <Button
                    id={`search-item-${index}`}
                    as="li"
                    aria-selected={selected ? true : undefined}
                    onMouseEnter={() => {
                      setActive(index);
                      eventRef.current = "mouse";
                    }}
                    onClick={() => {
                      // Can re-enable this for level-2 queries if we add it
                      // if (shouldCloseModal) {
                      el.action();
                      onClose();
                      // }
                    }}
                    // ref={menuNodes.ref(index)}
                    role="option"
                    key={el.name}
                    textAlign="left"
                    justifyContent="start"
                    sx={{
                      "display": "flex",
                      "alignItems": "center",
                      "minH": 7,
                      "mx": 2,
                      "px": 3,
                      "rounded": "md",
                      "bg": "none",
                      ".chakra-ui-da  rk &": { bg: "gray.600" },
                      "color": "gray.900",
                      "_hover": {
                        bg: "none",
                      },
                      "_selected": {
                        bg: "brandBlue.400",
                        color: "white",
                        mark: {
                          color: "white",
                          textDecoration: "underline",
                        },
                      },
                    }}
                  >
                    <el.icon mr={3} />
                    <Highlight query={[inputValue]} children={el.name} />
                    {el?.type === "favorite" && (
                      <Flex opacity=".5" ml="auto">
                        Favorited
                      </Flex>
                    )}
                  </Button>
                );
              })}
            </Box>
          </Flex>
        </ModalBody>
      </ModalContent>
    </Modal>
  );
};

export default EditActionModal;
