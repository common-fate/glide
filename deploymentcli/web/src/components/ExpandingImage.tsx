import {
  Image,
  ImageProps,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalOverlay,
  useDisclosure,
} from "@chakra-ui/react";
import React from "react";

type Props = ImageProps;

export const ExpandingImage: React.FC<Props> = (props) => {
  const { isOpen, onOpen, onClose } = useDisclosure();
  return (
    <>
      <Image {...props} onClick={onOpen} cursor="zoom-in" />
      <Modal isOpen={isOpen} onClose={onClose} size="7xl" isCentered>
        <ModalOverlay />
        <ModalContent>
          <ModalCloseButton />
          <ModalBody>
            <Image {...props} />
          </ModalBody>
        </ModalContent>
      </Modal>
    </>
  );
};
