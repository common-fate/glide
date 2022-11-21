import { DeleteIcon, EditIcon } from "@chakra-ui/icons";
import {
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  Avatar,
  Badge,
  Box,
  Button,
  ButtonGroup,
  Flex,
  HStack,
  IconButton,
  Input,
  InputGroup,
  InputRightElement,
  Link,
  Progress,
  Skeleton,
  SkeletonCircle,
  SkeletonText,
  Spinner,
  Stack,
  Text,
  Tooltip,
  useDisclosure,
  useToast,
  VStack,
  Wrap,
  WrapItem,
} from "@chakra-ui/react";
import axios from "axios";
import { intervalToDuration } from "date-fns";
import React, { createContext, useContext, useEffect, useState } from "react";
import ReactMarkdown from "react-markdown";
import { useUserListRequestsUpcoming } from "../utils/backend-client/default/default";
import {
  cancelRequest,
  reviewRequest,
  revokeRequest,
  useGetAccessInstructions,
  useGetAccessToken,
  useGetUser,
  useUserGetRequest,
} from "../utils/backend-client/end-user/end-user";
import {
  GrantStatus,
  RequestDetail,
  RequestStatus,
  ReviewDecision,
} from "../utils/backend-client/types";
import { Request } from "../utils/backend-client/types/request";
import { RequestTiming } from "../utils/backend-client/types/requestTiming";
import { useUser } from "../utils/context/userContext";
import { durationStringHoursMinutes } from "../utils/durationString";
import { renderTiming } from "../utils/renderTiming";
import { userName } from "../utils/userName";
import { CFReactMarkownCode } from "./CodeInstruction";
import { CopyableOption } from "./CopyableOption";
import { WarningIcon } from "./icons/Icons";
import { ProviderIcon } from "./icons/providerIcon";
import { InfoOption } from "./InfoOption";
import EditRequestTimeModal from "./modals/EditRequestTimeModal";
import RevokeConfirmationModal from "./modals/RevokeConfirmationModal";
import { StatusCell } from "./StatusCell";

interface RequestProps {
  request?: RequestDetail;
  isValidating: boolean;
  children?: React.ReactNode;
}

interface RequestContext {
  request: RequestDetail | undefined;
  isValidating: boolean;
  overrideTiming?: RequestTiming;
  setOverrideTiming: React.Dispatch<
    React.SetStateAction<RequestTiming | undefined>
  >;
}

const Context = createContext<RequestContext>({
  request: undefined,
  isValidating: true,
  setOverrideTiming: () => {
    undefined;
  },
});

export const RequestDisplay: React.FC<RequestProps> = ({
  request,
  isValidating,
  children,
}) => {
  const [overrideTiming, setOverrideTiming] = useState<RequestTiming>();

  useEffect(() => {
    // If its a schedule request in the past, set override timing default
    if (!overrideTiming) {
      if (
        request?.timing.startTime &&
        new Date(request.timing.startTime) < new Date()
      ) {
        setOverrideTiming({
          durationSeconds: request.timing.durationSeconds,
          startTime: new Date().toISOString(),
        });
      }
    }
  }, [request]);

  return (
    <Context.Provider
      value={{ overrideTiming, setOverrideTiming, request, isValidating }}
    >
      <Stack spacing={6} flex={1}>
        {children}
      </Stack>
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
  request: Request | RequestDetail | undefined,
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
  request: Request | RequestDetail | undefined;
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

export const RequestArgumentsDisplay: React.FC<{
  request: RequestDetail | undefined;
}> = ({ request }) => {
  if (request === undefined) {
    return <Skeleton minW="30ch" minH="6" mr="auto" />;
  }

  return (
    <VStack align={"left"}>
      <Text textStyle="Body/Medium">Request Details</Text>
      <Wrap>
        {Object.entries(request.arguments).map(([k, v]) => {
          return (
            <WrapItem>
              <VStack align={"left"}>
                <Text>{v.title}</Text>
                <InfoOption label={v.label} value={v.value} />
              </VStack>
            </WrapItem>
          );
        })}
      </Wrap>
    </VStack>
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
        label: `The access rule has been updated since this request was made. ${
          request?.canReview ? "You can still approve this request." : ""
        }`,
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
      <Stack spacing={2}>
        <Skeleton
          minW="30ch"
          minH="6"
          isLoaded={request?.accessRule !== undefined}
          mr="auto"
        >
          <HStack align="center" mr="auto">
            <ProviderIcon
              shortType={request?.accessRule.target.provider.type}
            />
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
        <RequestArgumentsDisplay request={request} />
        <Skeleton isLoaded={request !== undefined}>
          {/* @NOTE: this causes CLS, potential improvement */}
          {request?.reason && (
            <VStack align={"left"}>
              <Text textStyle="Body/Medium">Reason</Text>
              <Text
                color="neutrals.600"
                textStyle="Body/Small"
                data-testid="reason"
              >
                {request?.reason}
              </Text>
            </VStack>
          )}
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

  const [refreshInterval, setRefreshInterval] = useState(0);
  const { data: reqData } = useUserGetRequest(
    request?.grant != null ? request.id : "",

    {
      swr: {
        refreshInterval: refreshInterval,
      },
    }
  );

  useEffect(() => {
    if (reqData?.grant?.status == "PENDING") {
      if (
        reqData?.timing.startTime &&
        Date.parse(reqData.timing.startTime) > new Date().valueOf()
      ) {
        // This should make it refresh at least once just after its scheduled start time
        setRefreshInterval(
          Date.parse(reqData.timing.startTime) - new Date().valueOf() + 100
        );
      } else {
        setRefreshInterval(2000);
      }
    } else {
      setRefreshInterval(0);
    }
  }, [reqData]);

  if (!data || !data.instructions) {
    return null;
  }

  // Don't attempt to load a scheduled request until start time
  if (
    reqData?.timing.startTime &&
    Date.parse(reqData.timing.startTime) > new Date().valueOf()
  ) {
    return null;
  }

  if (reqData?.grant?.status === "PENDING") {
    return (
      <Stack>
        <Box textStyle="Body/Medium" id="access_instructions">
          Access Instructions
        </Box>
        <Text textStyle="Body/small" color="neutrals.600">
          Provisioning access
        </Text>
        <Progress size="xs" isIndeterminate hasStripe />
      </Stack>
    );
  }
  if (reqData?.grant?.status === "ACTIVE") {
    return (
      <Stack>
        <Box textStyle="Body/Medium" id="access_instructions">
          Access Instructions
        </Box>
        <ReactMarkdown
          components={{
            a: (props) => <Link target="_blank" rel="noreferrer" {...props} />,
            p: (props) => (
              <Text as="span" color="neutrals.600" textStyle={"Body/Small"}>
                {props.children}
              </Text>
            ),
            code: CFReactMarkownCode,
          }}
        >
          {data.instructions}
        </ReactMarkdown>
      </Stack>
    );
  }
  // Don't render anything
  return null;
};

export const RequestAccessToken: React.FC<{ reqId: string }> = ({ reqId }) => {
  const { data, error } = useGetAccessToken(reqId);

  const { data: reqData } = useUserGetRequest(reqId);

  const toast = useToast();
  // The Access tokne API returns a 404 for anyone other than the requestor or if there is no access token.
  // We treat the 404 as an indication tha there is no access token to display for this request
  if (error?.response?.status == 404) {
    return null;
  }
  if (!data) return <Spinner />;

  const handleClick = async () => {
    await navigator.clipboard.writeText(data);
    toast({
      title: "Access token copied to clipboard",
      status: "success",
      variant: "subtle",
      duration: 2200,
      isClosable: true,
    });
  };

  if (reqData?.grant?.status === "ACTIVE") {
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
  } else {
    return <></>;
  }
};

export const RequestTime: React.FC<{ canReview?: boolean }> = ({
  canReview,
}) => {
  const { request } = useContext(Context);

  return request ? (
    canReview ? (
      <_RequestOverridableTime />
    ) : (
      <_RequestTime />
    )
  ) : (
    // This let's us leverage the same Skeleton component
    <Stack spacing={2} minH="90px">
      <SkeletonText noOfLines={1} w="10ch" lineHeight="12px" />
      <SkeletonText noOfLines={1} w="10ch" lineHeight="12px" />
      <SkeletonText noOfLines={1} w="20ch" lineHeight="12px" />
    </Stack>
  );
};

export const _RequestTime: React.FC = () => {
  const { request } = useContext(Context);
  const timing = request?.timing;

  return (
    <Flex textStyle="Body/Small" flexDir="column" h="59px">
      <Box textStyle="Body/Medium">Duration</Box>
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
  );
};

/**
 * Similar to `<RequestTime />`, but allows the timing to be overridden during review.
 */
export const _RequestOverridableTime: React.FC = () => {
  const {
    request,
    setOverrideTiming,
    overrideTiming,
    isValidating,
  } = useContext(Context);
  const { onOpen, onClose, isOpen } = useDisclosure();
  const timing = request?.timing;

  const onUpdate = (timing: RequestTiming) => {
    setOverrideTiming(timing);
  };

  if (request?.status !== "PENDING") {
    return <_RequestTime />;
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
        <Text
          color="neutrals.600"
          textStyle="Body/Small"
          // noOfLines={1}
          p={1}
          pl={0}
          textDecoration={overrideTiming ? "line-through" : undefined}
        >
          {renderTiming(timing)}{" "}
        </Text>
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
      <EditRequestTimeModal
        handleSubmit={onUpdate}
        isOpen={isOpen}
        onClose={onClose}
        request={request}
      />
    </>
  );
};

export const RequestRequestor: React.FC = () => {
  const { request } = useContext(Context);
  const { data: requestor, isValidating } = useGetUser(
    request?.requestor ?? ""
  );

  /**
   * Load state requirements:
   * - show nothing if request?.requestor is undefined
   * - show skeleton if requestor is loading
   *
   * Improvements:
   * - reduce CLS by finding a way to absolute position component,
   * eliminating uncertainty of height
   */

  return !isValidating && requestor ? (
    <Flex textStyle="Body/Small" flexDir="column">
      <Box textStyle="Body/Medium">Requestor</Box>
      <Flex>
        <Avatar name={requestor.email} variant="withBorder" mr={2} size="xs" />
        <Text textStyle="Body/Small" mr={2}>
          {userName(requestor)}
        </Text>
        <Text
          color="neutrals.600"
          textStyle="Body/Small"
          maxW="20ch"
          noOfLines={1}
        >
          {requestor.email}
        </Text>
      </Flex>
    </Flex>
  ) : (
    <Flex flexDir="column">
      <SkeletonText noOfLines={1} w="10ch" height="18px" />
      <Flex alignItems="center">
        <SkeletonCircle size="24px" mr={2} />
        <SkeletonText noOfLines={1} w="10ch" lineHeight="12px" />
      </Flex>
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
  // @TODO: i can be replaced with a requestId and we can generate our own mutate hook
  onSubmitReview,
  canReview,
  focus,
}) => {
  const { request, overrideTiming, setOverrideTiming } = useContext(Context);
  const toast = useToast();
  const auth = useUser();
  const [isSubmitting, setIsSubmitting] = useState<ReviewDecision>();

  const onUpdate = (timing: RequestTiming) => {
    setOverrideTiming(timing);
  };

  const submitReview = async (decision: ReviewDecision) => {
    if (request === undefined) return;
    try {
      setIsSubmitting(decision);
      await reviewRequest(request.id, {
        decision,
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
      setIsSubmitting(undefined);
    }
  };

  // don't render any UI if we can't actually review the request.
  if (request?.status !== "PENDING") {
    return null;
  }

  const borderColor = focus !== undefined ? "brandGreen.300" : "neutrals.300";

  const { onOpen, onClose, isOpen } = useDisclosure();

  return (
    <Stack spacing={4}>
      <Text textStyle="Body/LargeBold">Review</Text>
      {
        // if the start time is in the past, show warning
        request.timing.startTime &&
          new Date() > new Date(request.timing.startTime) && (
            <Alert
              status="warning"
              rounded="md"
              borderColor={borderColor}
              borderRadius={"md"}
              borderWidth={"1px"}
              bg="white"
              alignItems="start"
              pb={8}
            >
              <WarningIcon boxSize="24px" mr={4} mt={1} />
              <Box>
                <AlertTitle mr={2} color="neutrals.800" fontWeight="medium">
                  This request is scheduled to start in the past
                </AlertTitle>
                <AlertDescription color="neutrals.600">
                  The scheduled start time for this request has already elapsed.{" "}
                  <strong style={{ fontWeight: "600" }}>Approving</strong> this
                  request will{" "}
                  <strong style={{ fontWeight: "600" }}>
                    activate access now
                  </strong>
                </AlertDescription>
                <Button
                  onClick={onOpen}
                  key={2}
                  rounded="full"
                  size="sm"
                  variant="outline"
                  position="absolute"
                  right={4}
                  bottom={4}
                  zIndex={2}
                  bg="white"
                >
                  Edit
                </Button>
              </Box>
            </Alert>
          )
      }
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
                    isLoading={isSubmitting === "APPROVED"}
                    isDisabled={isSubmitting === "DECLINED" || !canReview}
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
                    isDisabled={isSubmitting === "APPROVED"}
                    isLoading={isSubmitting === "DECLINED"}
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
      <EditRequestTimeModal
        handleSubmit={onUpdate}
        isOpen={isOpen}
        onClose={onClose}
        request={request}
      />
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
  // only display cancel if request is pending or is grant is still undefined
  if (
    (request?.status === "PENDING" && request.grant?.status == "PENDING") ||
    request?.grant == undefined
  ) {
    return (
      <ButtonGroup variant="outline" size="sm">
        {!request && <Skeleton rounded="full" w="64px" h="32px" />}
        <Button rounded="full" onClick={handleCancel}>
          Cancel
        </Button>
      </ButtonGroup>
    );
  } else {
    return null;
  }
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
