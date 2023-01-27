import { ArrowRightIcon } from "@chakra-ui/icons";
import { usePrefersReducedMotion, Box } from "@chakra-ui/react";
import { keyframes } from "@emotion/react";
import { motion } from "framer-motion";

const animationKeyframes = keyframes`
  0% { transform: translateX(0); opacity: 0.6; }
  100% { transform: translateX(0.2em); opacity: 1; }
`;

interface Props {
  animate?: boolean;
}

export const ConnectorArrow: React.FC<Props> = ({ animate }) => {
  const prefersReducedMotion = usePrefersReducedMotion();

  const animation =
    prefersReducedMotion || animate !== true
      ? `1s ease-in-out`
      : `${animationKeyframes} 1s ease-in-out infinite alternate`;

  return (
    <Box as={motion.div} animation={animation}>
      <ArrowRightIcon />
    </Box>
  );
};
