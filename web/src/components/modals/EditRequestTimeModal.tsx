import {
  Button,
  ButtonGroup,
  FormControl,
  FormHelperText,
  FormLabel,
  HStack,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  ModalProps,
  Stack,
  Text,
  useRadioGroup,
  UseRadioGroupProps,
} from "@chakra-ui/react";
import { format } from "date-fns";
import { useEffect, useMemo, useState } from "react";
import { Controller, useForm } from "react-hook-form";

import {
  Request,
  RequestAccessGroup,
  RequestAccessGroupTiming,
} from "../../utils/backend-client/types";

import { durationString } from "../../utils/durationString";
import { CFRadioBox } from "../CFRadioBox";
import { Days, DurationInput, Hours, Minutes, Weeks } from "../DurationInput";
export type When = "asap" | "scheduled";
export const WhenRadioGroup: React.FC<UseRadioGroupProps> = (props) => {
  const { getRootProps, getRadioProps } = useRadioGroup(props);
  const group = getRootProps();

  return (
    <HStack {...group}>
      <CFRadioBox {...getRadioProps({ value: "asap" })}>
        <Text textStyle="Body/Medium">ASAP</Text>
      </CFRadioBox>
      <CFRadioBox {...getRadioProps({ value: "scheduled" })}>
        <Text textStyle="Body/Medium">Scheduled</Text>
      </CFRadioBox>
    </HStack>
  );
};
type Props = {
  accessGroup: RequestAccessGroup;
  handleSubmit: (timing: RequestAccessGroupTiming) => void;
} & Omit<ModalProps, "children">;

interface ApproveRequestFormData {
  timing: RequestAccessGroupTiming;
  when: When;
}

const EditRequestTimeModal = ({ accessGroup, ...props }: Props) => {
  const [readableDuration, setReadableDuration] = useState<string>("1 hour");
  const methods = useForm<ApproveRequestFormData>();
  const when = methods.watch("when");
  const startTimeDate = methods.watch("timing.startTime");
  const now = useMemo(() => {
    const d = new Date();
    d.setSeconds(0, 0);
    return format(d, "yyyy-MM-dd'T'HH:mm");
  }, []);

  useEffect(() => {
    const data: ApproveRequestFormData = {
      timing: {
        durationSeconds: accessGroup.requestedTiming.durationSeconds,
        startTime: accessGroup.requestedTiming.startTime,
      },
      when: accessGroup.requestedTiming.startTime ? "scheduled" : "asap",
    };

    if (accessGroup.requestedTiming.startTime) {
      const d = new Date(Date.parse(accessGroup.requestedTiming.startTime));
      // This native datetime input needs a specific format as shown here, we take input in local time and it is converted to UTC for the api call
      data.timing.startTime = format(d, "yyyy-MM-dd'T'HH:mm");
    }
    methods.reset(data);
  }, []);

  const handleSubmit = async (data: ApproveRequestFormData) => {
    const startTime =
      data.when === "scheduled" && data.timing.startTime !== undefined
        ? new Date(data.timing.startTime).toISOString()
        : undefined;

    props.handleSubmit({
      durationSeconds: data.timing.durationSeconds,
      startTime,
    });
    props.onClose();
  };

  const maxDurationSeconds =
    accessGroup.accessRule.timeConstraints.maxDurationSeconds;
  return (
    <Modal {...props} size={maxDurationSeconds >= 3600 * 24 ? "xl" : "md"}>
      <form onSubmit={methods.handleSubmit(handleSubmit)}>
        <ModalOverlay />
        <ModalContent>
          <ModalHeader>Edit Request Time</ModalHeader>
          <ModalCloseButton />
          <ModalBody>
            <Stack>
              <FormControl pos="relative">
                <FormLabel textStyle="Body/Medium" fontWeight="normal">
                  Duration
                </FormLabel>
                <Controller
                  name="timing.durationSeconds"
                  control={methods.control}
                  rules={{
                    required: "Duration is required.",
                    max: maxDurationSeconds,
                    min: 60,
                  }}
                  render={({ field: { ref, ...rest } }) => {
                    return (
                      <DurationInput
                        {...rest}
                        max={maxDurationSeconds}
                        min={60}
                        defaultValue={
                          accessGroup.requestedTiming.durationSeconds
                        }
                        hideUnusedElements
                      >
                        <Weeks />
                        <Days />
                        <Hours />
                        <Minutes />
                        {maxDurationSeconds !== undefined && (
                          <Text textStyle={"Body/ExtraSmall"}>
                            Max {durationString(maxDurationSeconds)}
                            <br />
                            Min 1 min
                          </Text>
                        )}
                      </DurationInput>
                    );
                  }}
                />
                {/* <NumberInput
                  defaultValue={1}
                  min={0.01}
                  step={0.5}
                  max={12}
                  w="200px"
                  onChange={(s: string, n: number) => {
                    setReadableDuration(durationString(n * 3600));
                  }}
                >
                  <NumberInputField
                    bg="white"
                    {...methods.register("timing.durationSeconds")}
                  />
                  <NumberInputStepper>
                    <NumberIncrementStepper />
                    <NumberDecrementStepper />
                  </NumberInputStepper>
                </NumberInput>
                <FormHelperText color="neutrals.600">
                  {readableDuration}
                </FormHelperText> */}
              </FormControl>

              <FormControl
                pos="relative"
                id="when"
                isInvalid={methods.formState.errors.when !== undefined}
              >
                <FormLabel textStyle="Body/Medium" fontWeight="normal">
                  When
                </FormLabel>

                <Controller
                  name="when"
                  control={methods.control}
                  render={({ field }) => <WhenRadioGroup {...field} />}
                />
              </FormControl>

              {when === "scheduled" && (
                <FormControl>
                  <FormLabel textStyle="Body/Medium" fontWeight="normal">
                    Start Time
                  </FormLabel>

                  <Input
                    {...methods.register("timing.startTime")}
                    bg="white"
                    type="datetime-local"
                    min={now}
                    defaultValue={now}
                  />

                  {startTimeDate && (
                    <FormHelperText color="neutrals.600">
                      {new Date(startTimeDate).toString()}
                    </FormHelperText>
                  )}
                </FormControl>
              )}
            </Stack>
          </ModalBody>
          <ModalFooter minH={12}>
            <ButtonGroup rounded="full" spacing={2} ml="auto">
              <Button variant="outline" rounded="full" onClick={props.onClose}>
                Cancel
              </Button>
              <Button type="submit">Update</Button>
            </ButtonGroup>
          </ModalFooter>
        </ModalContent>
      </form>
    </Modal>
  );
};

export default EditRequestTimeModal;
