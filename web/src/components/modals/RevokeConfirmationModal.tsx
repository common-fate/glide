import {
  Button,
  ButtonGroup,
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
} from "@chakra-ui/react";
import { useState } from "react";

type Props = Omit<ModalProps, "children">;

interface _Props extends Props {
  action: string | undefined;
  onSubmit: () => Promise<void>;
}

const ConvertStatus = (status: string | undefined) => {
  switch (status) {
    case "ACTIVE":
      return "revoke";
      break;
    case "PENDING":
      return "cancel";
      break;
    default:
      break;
  }
};

const RevokeConfirmationModal = (props: _Props) => {
  const [loading, setLoading] = useState(false);

  return (
    <Modal {...props} isCentered>
      <ModalOverlay />
      <ModalContent>
        <ModalCloseButton />
        <ModalHeader>Revoke access</ModalHeader>
        <ModalBody>
          <Stack spacing="5">
            <Text>
              Are you sure you want to {ConvertStatus(props.action)} this grant?
            </Text>
          </Stack>
        </ModalBody>
        <ModalFooter minH={12}>
          <ButtonGroup spacing={2}>
            <Button
              isLoading={loading}
              onClick={() => {
                setLoading(true);
                props
                  .onSubmit()
                  .then(() => props.onClose())
                  .finally(() => setLoading(false));
              }}
              variant={"solid"}
              colorScheme="red"
              key={1}
              rounded="full"
            >
              Revoke
            </Button>
            <Button
              key={2}
              rounded="full"
              variant={"brandSecondary"}
              isDisabled={loading}
              onClick={props.onClose}
            >
              Cancel
            </Button>
          </ButtonGroup>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default RevokeConfirmationModal;
