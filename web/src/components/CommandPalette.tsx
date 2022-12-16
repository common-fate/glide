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
} from "@chakra-ui/react";
// import { useServiceKeys, useServiceMetaData } from "../../utils/apiHooks";
// import { ServiceMetadataPrivilege } from "../../iamzero-advisories/types";
import { useRef } from "react";
// import { useEditor } from "../../utils/context/EditorProvider";
import {
  AddIcon,
  ArrowForwardIcon,
  ArrowRightIcon,
  CheckIcon,
  DeleteIcon,
  EditIcon,
  PlusSquareIcon,
} from "@chakra-ui/icons";
import { useNavigate, useRouter } from "react-location";
// import QueryStringHighlight from "./QueryStringHighlight";

interface Props {}

type ICommand = {
  name: string;
  icon: ComponentWithAs<"svg", IconProps>;
  action: () => void;
  isAdminOnly?: boolean;
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
    {
      name: "Switch to admin",
      icon: ArrowForwardIcon,
      action: () => nav({ to: "/admin/" }),
      isAdminOnly: true,
    },
  ];

  /**
   Context aware
    - The palette knows where you are in the app
    - i.e. If you open the palette on an Access Request it displays a set of actions specific to that request e.g. 'Approve' or 'Revoke' request, or 'Copy access instructions'
   */

  const accessRequestActions: ICommand[] = [
    {
      name: "Approve request",
      icon: CheckIcon,
      action: () => undefined,
    },
    {
      name: "Revoke request",
      icon: DeleteIcon,
      action: () => undefined,
    },
  ];

  /**
   * This typing allows us to merge different search results,
   * specifying a search result type, key, description and the original node
   */
  type ResultsFormat<T, K> = {
    key: string;
    description: string;
    node: T;
    type: K;
  };

  //   type ServiceResult = ResultsFormat<ServiceMetadataPrivilege, "service">;
  //   type ServiceKeyResult = ResultsFormat<string, "serviceKey">;
  //   type CombinedResult = ServiceResult | ServiceKeyResult;

  //   const serviceKeyResults = React.useMemo(
  //     function getResults() {
  //       if (serviceKeys) {
  //         let serviceKeySearchResults: ServiceKeyResult[] = serviceKeys.map(
  //           (el) => ({
  //             key: el,
  //             description: el,
  //             node: el,
  //             type: "serviceKey",
  //           })
  //         );

  //         serviceKeySearchResults = matchSorter(
  //           serviceKeySearchResults,
  //           inputValue,
  //           {
  //             keys: ["node"],
  //             threshold: matchSorter.rankings.STARTS_WITH,
  //           }
  //         )
  //           .slice(0, 10)
  //           .filter((k) => k.key != selectedKey);

  //         return serviceKeySearchResults;
  //       } else return [];
  //     },
  //     [inputValue, serviceKeys]
  //   );

  //   const results = React.useMemo(
  //     function getResults() {
  //       if (serviceMetadata?.privileges) {
  //         let privilegeSearchResults: ServiceResult[] =
  //           serviceMetadata.privileges.map((el) => ({
  //             key: el.privilege,
  //             description: el.description,
  //             node: el,
  //             type: "service",
  //           }));

  //         // @TODO: turn this conditional check into a regex of whether the service is selected
  //         privilegeSearchResults = matchSorter(
  //           privilegeSearchResults,
  //           inputValue,
  //           {
  //             keys: [
  //               "node.description", // i.e. 'Grants permission to...' - the long description
  //               "node.access_level", // i.e. 'Write' - used to filter by category
  //               (key) => selectedKey + ":" + key.node.privilege, // i.e. 's3:' + 'AllocateAddress' - the privileges key
  //             ],
  //             threshold: matchSorter.rankings.CONTAINS,
  //           }
  //         ).slice(0, 50);

  //         return privilegeSearchResults;
  //       } else return [];
  //     },
  //     [inputValue, serviceMetadata]
  //   );

  //   const loading = isValidating && !serviceMetadata;
  const loading = false;

  return (
    <Modal isOpen={isOpen} onClose={onClose} size="md">
      <ModalOverlay />
      <ModalContent overflow="hidden">
        {/* <ModalCloseButton zIndex={999} size="sm" /> */}
        {/* <ModalHeader fontSize="md" pb={2}>
          Add an action
        </ModalHeader> */}
        <ModalBody p={0} position="relative" pb={3} h="100%" maxH="80vh">
          <Flex flex={1} position="relative" flexDir="column" pb={4}>
            <InputGroup>
              <Input
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
              {[].length > 0 && (
                <Box fontSize="sm" flex="1 0 20%" mt={5} px={6}>
                  <Text opacity={0.52} display="inline">
                    Filter by Service:
                  </Text>
                  <Kbd float="right" opacity=".4">
                    Tab
                  </Kbd>
                  <br />

                  <HStack
                    spacing={1}
                    display="inline"
                    ml={1}
                    opacity={1}
                    position="relative"
                  >
                    {[...[]]
                      //   .filter((result) => result.type == "serviceKey")
                      .slice(0, 5)
                      .map((el, i) => (
                        <Button
                          opacity={0.52}
                          _selected={{
                            opacity: "1 !important",
                          }}
                          _focus={{
                            opacity: "1 !important",
                            boxShadow: "outline",
                          }}
                          // @TODO: we may want to improve these selective stylings a bit more
                          // opacity={
                          // 	i == 0
                          // 		? '1 !important'
                          // 		: '.4'
                          // }
                          // colorScheme={
                          // 	i == 0
                          // 		? 'cyan'
                          // 		: 'gray'
                          // }
                          size="xs"
                          //   key={i}
                          //   onClick={() => {
                          // setSelectedKey(el.node);
                          // setInputValue(el.node + ":");
                          // inputRef?.current?.focus();
                          //   }}
                        >
                          {/* {el.node} */}
                        </Button>
                      ))}
                    <chakra.span opacity={0.52} pos="absolute">
                      {[].length > 5 && "..."}
                    </chakra.span>
                  </HStack>
                </Box>
              )}

              <Box fontSize="sm" mt={5} px={6}>
                {/* <Text opacity={0.52} display="inline">
                  Showing suggestions for&nbsp;
                  <Badge colorScheme="cyan">{selectedKey}</Badge>
                </Text> */}
              </Box>

              {[...otherActions, ...accessRequestActions]
                /**
                 * @NOTE: if there's a performant way to adding sorting
                 * based on privilege Type > type results.length
                 * then we should implement it
                 * */
                // .sort((a, b) =>
                // 	results.filter(
                // 		(x) => x.node?.access_level == a
                // 	).length >
                // 	results?.filter(
                // 		(x) => x.node?.access_level == b
                // 	).length
                // 		? 1
                // 		: 0
                // )
                .map((item) => {
                  // Total count
                  //   let [] = serviceMetadata?.privileges.filter(
                  //     (a) => a.access_level == privilegeType
                  //   );
                  //   // Filtered search results
                  //   let [] = results.filter(
                  //     (p) =>
                  //       p.type == "service" &&
                  //       p.node.access_level == privilegeType
                  //   );

                  return (
                    <Box mt={5} key={item.name} mx={6}>
                      <Text
                        fontSize="sm"
                        opacity={0.7}
                        display="flex"
                        justifyContent="space-between"
                      >
                        <span>
                          <Highlight
                            query={[inputValue]}
                            children={item.name}
                          />
                        </span>
                        <span>
                          {/* add in recent/bookmarked subtext here */}
                        </span>
                      </Text>
                    </Box>
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
