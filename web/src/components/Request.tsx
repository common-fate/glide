import { CheckIcon, CopyIcon, DeleteIcon, EditIcon } from "@chakra-ui/icons";
import {
  Avatar,
  Badge,
  Box,
  Button,
  ButtonGroup,
  Code,
  Flex,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Link,
  Skeleton,
  SkeletonText,
  Spinner,
  Spacer,
  Stack,
  Text,
  Tooltip,
  useClipboard,
  useDisclosure,
  useToast,
} from "@chakra-ui/react";
import axios from "axios";
import { intervalToDuration } from "date-fns";
import React, { createContext, useContext, useState } from "react";
import ReactDom from "react-dom";
import ReactMarkdown from "react-markdown";
import { durationStringHoursMinutes } from "../utils/durationString";
import { useUserListRequestsUpcoming } from "../utils/backend-client/default/default";
import {
  cancelRequest,
  reviewRequest,
  revokeRequest,
  useGetAccessInstructions,
  useGetAccessToken,
  useGetUser,
} from "../utils/backend-client/end-user/end-user";
import {
  GrantStatus,
  RequestDetail,
  RequestStatus,
  ReviewDecision,
} from "../utils/backend-client/types";
import { RequestTiming } from "../utils/backend-client/types/requestTiming";
import { Request } from "../utils/backend-client/types/request";
import { useUser } from "../utils/context/userContext";
import { renderTiming } from "../utils/renderTiming";
import { userName } from "../utils/userName";
import { ProviderIcon } from "./icons/providerIcon";
import EditRequestTimeModal from "./modals/EditRequestTimeModal";
import RevokeConfirmationModal from "./modals/RevokeConfirmationModal";
import { RequestStatusCell, StatusCell } from "./StatusCell";
import rehypeRaw from "rehype-raw";
import { CodeProps } from "react-markdown/lib/ast-to-react";

interface RequestProps {
  request?: RequestDetail;
  children?: React.ReactNode;
}

interface RequestContext {
  request: RequestDetail | undefined;
  overrideTiming?: RequestTiming;
  setOverrideTiming: React.Dispatch<
    React.SetStateAction<RequestTiming | undefined>
  >;
}

const Context = createContext<RequestContext>({
  request: undefined,
  setOverrideTiming: () => {
    undefined;
  },
});

export const RequestDisplay: React.FC<RequestProps> = ({
  request,
  children,
}) => {
  const [overrideTiming, setOverrideTiming] = useState<RequestTiming>();

  return (
    <Context.Provider value={{ overrideTiming, setOverrideTiming, request }}>
      <Stack spacing={6}>{children}</Stack>
    </Context.Provider>
  );
};

interface RequestDetailProps {
  children?: React.ReactNode;
}

interface VersionDisplay {
  badge: string;
  label: string;
}

const getStatus = (
  request: Request | undefined,
  activeTimeString: string | undefined
) => {
  if (activeTimeString !== undefined) {
    return activeTimeString;
  }
  if (request?.grant !== undefined) {
    if (request?.grant.status === "PENDING") {
      return request?.status;
    }

    return request.grant.status;
  }
  if (
    request?.status === "APPROVED" &&
    request.approvalMethod === "AUTOMATIC"
  ) {
    return "Automatically approved";
  }

  return request?.status;
};

export const RequestStatusDisplay: React.FC<{
  request: Request | undefined;
}> = ({ request }) => {
  const activeTimeString =
    request?.grant && request?.grant.status === "ACTIVE"
      ? "Active for the next " +
        durationStringHoursMinutes(
          intervalToDuration({
            start: new Date(),
            end: new Date(Date.parse(request.grant.end)),
          })
        )
      : undefined;

  const status = getStatus(request, activeTimeString);

  return (
    <StatusCell
      value={status}
      success={[
        activeTimeString ?? "",
        GrantStatus.ACTIVE,
        "Automatically approved",
      ]}
      info={[RequestStatus.APPROVED]}
      danger={[
        RequestStatus.DECLINED,
        RequestStatus.CANCELLED,
        GrantStatus.REVOKED,
      ]}
      warning={RequestStatus.PENDING}
      textStyle="Body/Small"
    />
  );
};

export const RequestDetails: React.FC<RequestDetailProps> = ({ children }) => {
  const { request } = useContext(Context);

  const version: VersionDisplay = request?.accessRule.isCurrent
    ? {
        badge: "Latest",
        label:
          "This request was made for the current version of this access rule.",
      }
    : {
        badge: "Old Version",
        label:
          "The access rule has been updated since this request was made. You can still approve this request.",
      };

  return (
    <Stack
      rounded="md"
      bg="neutrals.100"
      flexDir="column"
      w="100%"
      // w={{ base: "100%", md: "500px", lg: "716px" }}
      p={8}
      spacing={6}
    >
      <RequestStatusDisplay request={request} />
      <Stack spacing={1}>
        <Skeleton
          minW="30ch"
          minH="6"
          isLoaded={request?.accessRule !== undefined}
          mr="auto"
        >
          <HStack align="center" mr="auto">
            <ProviderIcon provider={request?.accessRule.target.provider} />
            <Text textStyle="Body/LargeBold">{request?.accessRule?.name}</Text>
            <Tooltip label={version.label}>
              <Badge
                fontSize={"9px"}
                fontWeight="normal"
                variant={"outline"}
                // ensure this remains focusable for accessibility, as the tooltip has the full context on what the version actually means.
                tabIndex={0}
                cursor="default"
              >
                {version.badge}
              </Badge>
            </Tooltip>
          </HStack>
        </Skeleton>
        <Skeleton isLoaded={request !== undefined}>
          <Flex
            color="neutrals.600"
            textStyle="Body/Medium"
            data-testid="reason"
          >
            {request?.reason}
          </Flex>
        </Skeleton>
      </Stack>
      {children}
    </Stack>
  );
};

export const RequestAccessInstructions: React.FC = () => {
  const { request } = useContext(Context);
  const { data } = useGetAccessInstructions(
    request?.grant != null ? request.id : ""
  );

  if (!data || !data.instructions) {
    return null;
  }

  return (
    <Stack>
      <Box textStyle="Body/Medium" mb={2}>
        Access Instructions
      </Box>
      return (
      <ReactMarkdown
        rehypePlugins={[rehypeRaw]}
        // remarkPlugins={[remarkGfm]}
        skipHtml={false}
        components={{
          a: (props) => <Link target="_blank" rel="noreferrer" {...props} />,
          p: (props) => (
            <Text as="span" color="neutrals.600" textStyle={"Body/Small"}>
              {props.children}
            </Text>
          ),
          code: CodeInstruction,
        }}
        children={data.instructions}
      />
    </Stack>
  );
};

export const RequestAccessToken = () => {
  const { request } = useContext(Context);
  const { data } = useGetAccessToken(request?.id);

  // const [token, setToken] = useState<string>();

  const toast = useToast();

  if (!data) return <Spinner />;

  const handleClick = () => {
    navigator.clipboard.writeText(data);
    toast({
      title: "Access token copied to clipboard",
      status: "success",
      variant: "subtle",
      duration: 2200,
      isClosable: true,
    });
  };

  return (
    <Stack>
      <Box textStyle="Body/Medium" mb={2}>
        Access Token
      </Box>
      <InputGroup size="md" bg="white" maxW="400px">
        <Input pr="4.5rem" type={"password"} value={data} readOnly />
        <InputRightElement width="4.5rem" pr={1}>
          <Button h="1.75rem" size="sm" onClick={handleClick}>
            {"Copy"}
          </Button>
        </InputRightElement>
      </InputGroup>
    </Stack>
  );
};

const CodeInstruction: React.FC<CodeProps> = (props) => {
  const { children, node } = props;
  let value = "";
  if (node.children.length == 1 && node.children[0].type == "text") {
    value = node.children[0].value;
  }

  const { hasCopied, onCopy } = useClipboard(value);
  return (
    <Stack>
      <Code
        padding={0}
        bg="white"
        borderRadius="8px"
        borderColor="neutrals.300"
        borderWidth="1px"
      >
        <Flex
          borderColor="neutrals.300"
          borderBottomWidth="1px"
          py="8px"
          px="16px"
          minH="36px"
        >
          <Spacer />
          <IconButton
            variant="ghost"
            h="20px"
            icon={hasCopied ? <CheckIcon /> : <CopyIcon />}
            onClick={onCopy}
            aria-label={"Copy"}
          />
        </Flex>
        <Text
          overflowX="auto"
          color="neutrals.700"
          padding={4}
          whiteSpace="pre-wrap"
        >
          {children}
        </Text>
      </Code>
    </Stack>
  );
};
export const RequestTime: React.FC = () => {
  const { request } = useContext(Context);
  const timing = request?.timing;

  return request ? (
    <Flex textStyle="Body/Small" flexDir="column" h="59px">
      <Box textStyle="Body/Medium" mb={2}>
        Duration
      </Box>
      <Text
        color="neutrals.600"
        textStyle="Body/Small"
        noOfLines={1}
        p={1}
        pl={0}
      >
        {renderTiming(timing)}{" "}
      </Text>
    </Flex>
  ) : (
    <Stack spacing={2}>
      <SkeletonText noOfLines={1} w="10ch" />
      <SkeletonText noOfLines={1} w="20ch" />
    </Stack>
  );
};

/**
 * Similar to `<RequestTime />`, but allows the timing to be overridden during review.
 */
export const RequestOverridableTime: React.FC = () => {
  const { request, setOverrideTiming, overrideTiming } = useContext(Context);
  const { onOpen, onClose, isOpen } = useDisclosure();
  const timing = request?.timing;

  const onUpdate = (timing: RequestTiming) => {
    setOverrideTiming(timing);
  };

  if (request?.status !== "PENDING") {
    return <RequestTime />;
  }

  return (
    <>
      <Flex textStyle="Body/Small" flexDir="column">
        <Box textStyle="Body/Medium" mb={1} pos="relative">
          Timing
          <IconButton
            onClick={onOpen}
            variant={"ghost"}
            icon={<EditIcon />}
            aria-label="edit"
            size="xs"
            pos="absolute"
            left={14}
            rounded="full"
          />
        </Box>
        <Skeleton isLoaded={request !== undefined}>
          <Text
            color="neutrals.600"
            textStyle="Body/Small"
            // noOfLines={1}
            p={1}
            pl={0}
            textDecoration={
              overrideTiming !== undefined ? "line-through" : undefined
            }
          >
            {renderTiming(timing)}{" "}
          </Text>
        </Skeleton>
        {overrideTiming && (
          <Text
            color="neutrals.600"
            textStyle="Body/Small"
            noOfLines={1}
            p={1}
            pl={0}
            fontStyle={"italic"}
          >
            {renderTiming(overrideTiming)}{" "}
            <IconButton
              onClick={() => setOverrideTiming(undefined)}
              variant={"ghost"}
              icon={<DeleteIcon />}
              aria-label="undo"
              size="xs"
            />
          </Text>
        )}
      </Flex>
      {request && (
        <EditRequestTimeModal
          handleSubmit={onUpdate}
          isOpen={isOpen}
          onClose={onClose}
          request={request}
        />
      )}
    </>
  );
};

export const RequestRequestor: React.FC = () => {
  const { request } = useContext(Context);
  const { data: requestor } = useGetUser(request?.requestor ?? "");

  return (
    <Flex textStyle="Body/Small" flexDir="column">
      <Box textStyle="Body/Medium" mb={2}>
        Requestor
      </Box>
      <Skeleton w="30ch" isLoaded={requestor !== undefined}>
        {requestor && (
          <Flex>
            <Avatar
              name={requestor.email}
              variant="withBorder"
              mr={2}
              size="xs"
            />
            <Text textStyle="Body/Small">{userName(requestor)}</Text>
            <Text
              color="neutrals.600"
              textStyle="Body/Small"
              maxW="20ch"
              noOfLines={1}
            >
              {requestor.email}
            </Text>
          </Flex>
        )}
      </Skeleton>
    </Flex>
  );
};

interface ReviewButtonsProps {
  canReview: boolean;
  onSubmitReview?: () => void;
  /** if set, autofocuses to the approve or close review button and displays 'Confirm your review'. */
  focus?: "approve" | "close";
}

export const RequestReview: React.FC<ReviewButtonsProps> = ({
  onSubmitReview,
  canReview,
  focus,
}) => {
  const { request, overrideTiming } = useContext(Context);
  const toast = useToast();
  const auth = useUser();
  const [comment, setComment] = useState<string>("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  const submitReview = async (decision: ReviewDecision) => {
    if (request === undefined) return;
    try {
      setIsSubmitting(true);
      await reviewRequest(request.id, {
        decision,
        comment,
        overrideTiming: overrideTiming,
      });
      toast({
        title: decision === "APPROVED" ? "Request approved" : "Request closed",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      onSubmitReview && onSubmitReview();
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }

      toast({
        title: "Error submitting review",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  // don't render any UI if we can't actually review the request.
  if (request?.status !== "PENDING") {
    return null;
  }

  const borderColor = focus !== undefined ? "brandGreen.300" : "neutrals.300";

  return (
    <Stack spacing={4}>
      <Text textStyle="Body/LargeBold">Review</Text>
      <HStack spacing={3}>
        <Avatar
          variant="withBorder"
          size="sm"
          name={auth?.user?.email}
          alignSelf="flex-start"
          mt={1}
        />
        <Stack
          flexGrow={1}
          spacing={4}
          p={4}
          borderColor={borderColor}
          borderRadius={"md"}
          borderWidth={"1px"}
          position={"relative"}
          _before={{
            position: "absolute",
            height: "16px",
            width: "8px",
            top: "11px",
            left: "-8px",
            right: "100%",
            backgroundColor: borderColor,
            content: `" "`,
            clipPath: "polygon(0 50%, 100% 0, 100% 100%)",
          }}
          _after={{
            ml: "1.5px",
            position: "absolute",
            height: "16px",
            width: "8px",
            top: "11px",
            left: "-8px",
            right: "100%",
            backgroundColor: "white",
            content: `" "`,
            clipPath: "polygon(0 50%, 100% 0, 100% 100%)",
          }}
        >
          {/* <Input
            isDisabled={!canReview}
            borderColor={"neutrals.300"}
            value={comment}
            onChange={(e: {
              target: { value: React.SetStateAction<string> };
            }) => setComment(e.target.value)}
            placeholder="Leave an optional review comment"
            bg="white"
          />*/}
          <Tooltip
            isDisabled={canReview}
            label="Users cannot review their own requests."
          >
            <ButtonGroup
              variant="outline"
              size="sm"
              alignSelf={"flex-end"}
              isDisabled={!canReview}
            >
              {request?.status === "PENDING" && (
                <>
                  {focus !== undefined && (
                    <Text mt={2} textStyle="Body/ExtraSmall">
                      Confirm your review
                    </Text>
                  )}
                  <Button
                    data-testid="approve"
                    isLoading={isSubmitting}
                    autoFocus={focus === "approve"}
                    variant={"brandPrimary"}
                    key={1}
                    rounded="full"
                    onClick={() => submitReview("APPROVED")}
                  >
                    Approve
                  </Button>
                  <Button
                    data-testid="decline"
                    isLoading={isSubmitting}
                    autoFocus={focus === "close"}
                    key={2}
                    rounded="full"
                    onClick={() => submitReview("DECLINED")}
                  >
                    Close Request
                  </Button>
                </>
              )}
            </ButtonGroup>
          </Tooltip>
        </Stack>
      </HStack>
    </Stack>
  );
};

export const RequestCancelButton: React.FC = () => {
  const { request } = useContext(Context);
  const toast = useToast();
  const { mutate } = useUserListRequestsUpcoming();

  const handleCancel = async () => {
    if (request === undefined) return;
    try {
      await cancelRequest(request.id, {});
      void mutate();
      toast({
        title: "Request cancelled",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }

      toast({
        title: "Error cancelling request",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };

  return (
    <ButtonGroup variant="outline" size="sm">
      {!request && <Skeleton rounded="full" w="64px" h="32px" />}
      {request?.status === "PENDING" && (
        <Button rounded="full" onClick={handleCancel}>
          Cancel
        </Button>
      )}
    </ButtonGroup>
  );
};

interface RevokeButtonsProps {
  onSubmitRevoke?: () => void;
}

export const RequestRevoke: React.FC<RevokeButtonsProps> = ({
  onSubmitRevoke,
}) => {
  const { request } = useContext(Context);

  const toast = useToast();

  const submitRevoke = async () => {
    if (request === undefined) return;
    try {
      await revokeRequest(request.id, {});
      toast({
        title: "Deactivated grant",
        status: "success",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
      onSubmitRevoke && onSubmitRevoke();
    } catch (err) {
      let description: string | undefined;
      if (axios.isAxiosError(err)) {
        // @ts-ignore
        description = err?.response?.data.error;
      }

      toast({
        title: "Error deactivating grant",
        description,
        status: "error",
        variant: "subtle",
        duration: 2200,
        isClosable: true,
      });
    }
  };

  const revokeConfirmationDisclosure = useDisclosure();

  const renderButton = (status: string, grantId: string) => {
    switch (status) {
      case "ACTIVE":
        return (
          <ButtonGroup variant="outline" size="sm" alignSelf={"flex-end"}>
            <>
              <Button
                key={2}
                rounded="full"
                onClick={() => revokeConfirmationDisclosure.onOpen()}
              >
                Revoke
              </Button>
            </>
          </ButtonGroup>
        );
      case "PENDING":
        return (
          <ButtonGroup variant="outline" size="sm" alignSelf={"flex-end"}>
            <>
              <Button
                key={2}
                rounded="full"
                onClick={() => revokeConfirmationDisclosure.onOpen()}
              >
                Cancel
              </Button>
            </>
          </ButtonGroup>
        );
    }
    return null;
  };

  // don't render any UI if we can't actually revoke the request.
  const canRevoke =
    request?.grant?.status == "ACTIVE" || request?.grant?.status == "PENDING";
  if (request?.status !== "APPROVED" || !canRevoke) {
    return null;
  }

  return (
    <>
      <Stack spacing={4}>
        <HStack spacing={3}>
          {request?.grant && renderButton(request?.grant?.status, request.id)}
        </HStack>
      </Stack>
      <RevokeConfirmationModal
        onSubmit={submitRevoke}
        isOpen={revokeConfirmationDisclosure.isOpen}
        onClose={revokeConfirmationDisclosure.onClose}
        action={request.grant?.status}
      />
    </>
  );
};
