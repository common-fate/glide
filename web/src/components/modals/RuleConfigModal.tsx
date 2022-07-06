import {
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  ModalProps,
} from "@chakra-ui/react";
import { AccessRule } from "../../utils/backend-client/types";

type Props = { rule?: AccessRule } & Omit<ModalProps, "children">;

const RuleConfigModal = ({ rule, ...props }: Props) => {
  return (
    <Modal {...props}>
      <ModalOverlay />
      <ModalContent>
        <ModalCloseButton />
        <ModalHeader mt={10}>Rule Config</ModalHeader>
        <ModalBody>{JSON.stringify(rule)?.slice(0, 400)}</ModalBody>
        <ModalFooter minH={12}>
          {/* <Button mr={3} isLoading={loading}>
            Create user
          </Button> */}
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
};

export default RuleConfigModal;
