import { ArrowBackIcon } from "@chakra-ui/icons";
import { Button, Container, IconButton, Stack, Text } from "@chakra-ui/react";
import { Link } from "react-location";
import ReactMarkdown from "react-markdown";
import { CodeInstruction } from "../../../../components/CodeInstruction";
import { AdminLayout } from "../../../../components/Layout";

const Page = () => {
  return (
    <AdminLayout>
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
          to="/admin/providers"
        />
        <Text as="h4" textStyle="Heading/H4">
          Provider setup complete
        </Text>
      </Stack>
      <Container
        my={12}
        // This prevents unbounded widths for small screen widths
        minW={{ base: "100%", xl: "container.xl" }}
        overflowX="auto"
      >
        <Stack
          px={8}
          py={8}
          bg="neutrals.100"
          rounded="md"
          w="100%"
          spacing={8}
        >
          <Text>
            To apply the changes, update your deployment of Granted by running:
          </Text>
          <ReactMarkdown
            components={{
              a: (props) => (
                <Link target="_blank" rel="noreferrer" {...props} />
              ),
              p: (props) => (
                <Text as="span" color="neutrals.600" textStyle={"Body/Small"}>
                  {props.children}
                </Text>
              ),
              code: CodeInstruction as any,
            }}
          >
            ``` gdeploy update ```
          </ReactMarkdown>
          <Button as={Link} to="/admin/providers">
            Back to Providers
          </Button>
        </Stack>
      </Container>
    </AdminLayout>
  );
};

export default Page;
