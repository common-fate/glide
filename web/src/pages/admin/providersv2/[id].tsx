import { ArrowBackIcon } from "@chakra-ui/icons";
import {
  Box,
  Center,
  Circle,
  Code,
  Container,
  Flex,
  IconButton,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Skeleton,
  Text,
  VStack,
  useDisclosure,
} from "@chakra-ui/react";
import { useMemo, useState } from "react";
import { Helmet } from "react-helmet";
import { Link, useMatch } from "react-location";
import { Column } from "react-table";
import { AdminLayout } from "../../../components/Layout";
import { TableRenderer } from "../../../components/tables/TableRenderer";
import { useAdminGetTargetGroup } from "../../../utils/backend-client/admin/admin";
import { useAdminListTargetRoutes } from "../../../utils/backend-client/default/default";
import { Diagnostic, TargetRoute } from "../../../utils/backend-client/types";
import React from "react";
import { targetGroupFromToString } from ".";

const AdminRoutesTable = () => {
  const {
    params: { id: groupId },
  } = useMatch();
  const { data } = useAdminListTargetRoutes(groupId, {});

  const diagnosticModal = useDisclosure();
  const [diagnosticText, setDiagnosticText] = useState("");

  const cols: Column<TargetRoute>[] = useMemo(
    () => [
      {
        accessor: "handlerId",
        Header: "Handler",
      },
      {
        accessor: "kind",
        Header: "Kind",
      },
      {
        accessor: "priority",
        Header: "Priority",
      },
      {
        accessor: "valid",
        Header: "Status",
        Cell: ({ value }) => (
          <Flex minW="75px" align="center">
            <Circle
              bg={value ? "actionSuccess.200" : "actionWarning.200"}
              size="8px"
              mr={2}
            />
            <Text as="span">{value ? "Valid" : "Invalid"}</Text>
          </Flex>
        ),
      },

      {
        accessor: "diagnostics",
        Header: "Diagnostics",
        Cell: ({ value }) => {
          // Strip out the code from the diagnostics, it's currently an empty field
          const strippedCode = JSON.stringify(
            (value as Partial<Diagnostic>[]).map((v) => {
              delete v["code"];
              return v;
            })
          );

          // Forces the modal prompt to always show
          const expandCode = strippedCode.length !== 2;

          const handleClick = () => {
            if (expandCode) {
              diagnosticModal.onOpen();
              setDiagnosticText(strippedCode);
            }
          };

          return (
            <Code
              maxW="400px"
              rounded="md"
              fontSize="sm"
              userSelect={expandCode ? "none" : "auto"}
              p={2}
              noOfLines={3}
              onClick={handleClick}
              position="relative"
              _hover={{
                "backgroundColor": expandCode ? "gray.600" : "gray.200",
                "cursor": expandCode ? "pointer" : "default",
                "#expandCode": {
                  display: "block",
                },
              }}
            >
              {expandCode && (
                <Text
                  id="expandCode"
                  display="none"
                  position="absolute"
                  left="50%"
                  top="50%"
                  transform="translate(-50%, -50%)"
                  zIndex={2}
                  size="md"
                  color="gray.50"
                >
                  Expand code
                </Text>
              )}
              {strippedCode}
            </Code>
          );
        },
      },
    ],
    []
  );

  return TableRenderer<TargetRoute>({
    columns: cols,
    data: data?.routes,
    emptyText: "No Routes have been set up yet.",
    linkTo: false,
    additionalChildren: (
      <Modal isOpen={diagnosticModal.isOpen} onClose={diagnosticModal.onClose}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Diagnostics</ModalHeader>
          <ModalCloseButton />
          <ModalBody pb={4}>
            <Code rounded="md" minH="200px" fontSize="sm" p={2}>
              {diagnosticText}
            </Code>
          </ModalBody>
        </ModalContent>
      </Modal>
    ),
  });
};

const Index = () => {
  const {
    params: { id: groupId },
  } = useMatch();
  // const ruleId = typeof query?.id == "string" ? query.id : "";
  const { data, isValidating, error } = useAdminGetTargetGroup(groupId);

  return (
    <>
      <AdminLayout>
        <Helmet>
          <title>{data?.id}</title>
        </Helmet>
        <Center borderBottom="1px solid" borderColor="neutrals.200" h="80px">
          <IconButton
            as={Link}
            to={"/admin/providersv2"}
            aria-label="Go back"
            pos="absolute"
            left={4}
            icon={<ArrowBackIcon />}
            rounded="full"
            variant="ghost"
          />

          <Text as="h4" textStyle="Heading/H4">
            {data?.id}
          </Text>
        </Center>

        <Container maxW="container.xl" py={16}>
          <Center>
            <Flex
              direction={["column", "row"]}
              rounded="md"
              bg="neutrals.100"
              p={8}
            >
              <VStack align={"left"} spacing={3} flex={1} mr={4}>
                <Text textStyle="Body/Medium">Name</Text>
                <Text textStyle="Body/Small">{data?.id}</Text>
                <Text textStyle="Body/Medium">Target Schema</Text>
                <Text textStyle="Body/Small">
                  {data?.from ? targetGroupFromToString(data.from) : ""}
                </Text>
                <Text textStyle="Body/Medium">Routes</Text>

                <AdminRoutesTable />
              </VStack>

              <Box mb={24}>
                {!data && (
                  <Container pt={12} maxW="container.md">
                    <VStack w="100%" spacing={6}>
                      <Skeleton h="150px" w="100%" rounded="md" />
                      <Skeleton h="150px" w="100%" rounded="md" />
                      <Skeleton h="150px" w="100%" rounded="md" />
                      <Skeleton h="150px" w="100%" rounded="md" />
                      <Skeleton h="150px" w="100%" rounded="md" />
                    </VStack>
                  </Container>
                )}
              </Box>
            </Flex>
          </Center>
        </Container>
      </AdminLayout>
    </>
  );
};

export default Index;
