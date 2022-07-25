import {
  Button,
  ButtonGroup,
  FormControl,
  FormHelperText,
  FormLabel,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  ModalProps,
  NumberDecrementStepper,
  NumberIncrementStepper,
  NumberInput,
  NumberInputField,
  NumberInputStepper,
  Stack,
  Text,
} from "@chakra-ui/react";
import { format } from "date-fns";
import { useEffect, useMemo, useState } from "react";
import { Controller, useForm } from "react-hook-form";
import { When, WhenRadioGroup } from "../../pages/access/request/[id]";
import { Request, RequestDetail } from "../../utils/backend-client/types";
import { RequestTiming } from "../../utils/backend-client/types/requestTiming";
import { durationString } from "../../utils/durationString";
import { DurationInput, Hours, Minutes } from "../DurationInput";
type Props = {
  request: RequestDetail;
  handleSubmit: (timing: RequestTiming) => void;
} & Omit<ModalProps, "children">;

interface ApproveRequestFormData {
  timing: RequestTiming;
  when: When;
}

const EditRequestTimeModal = ({ request, ...props }: Props) => {
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
    let data: ApproveRequestFormData = {
      timing: {
        durationSeconds: request.timing.durationSeconds,
        startTime: request.timing.startTime,
      },
      when: request.timing.startTime ? "scheduled" : "asap",
    };

    if (request.timing.startTime) {
      const d = new Date(Date.parse(request.timing.startTime));
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
    request.accessRule.timeConstraints.maxDurationSeconds;
  return (
    <Modal {...props}>
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
                        defaultValue={request.timing.durationSeconds}
                      >
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
