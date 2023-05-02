import { Box, BoxProps, createIcon } from "@chakra-ui/react";

// blue, pink, green shapes
export const BlurShapes = {
  Blue: createIcon({
    viewBox: "0 0 476 426",
    path: (
      <svg width="476" height="426" viewBox="0 0 476 426" fill="none">
        <path
          d="M476 97.7605L0 0L193.949 138.003L63.6298 426L367.095 270.925L452.139 367.669L476 97.7605Z"
          fill="#2E7FFF"
        />
      </svg>
    ),
  }),
  Pink: createIcon({
    viewBox: "0 0 372 310",
    path: (
      <svg
        width="372"
        height="310"
        viewBox="0 0 372 310"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M223.529 0L91.5186 20.0457L0 169.68L174.478 310L296.942 211.796L174.478 140.523L372 128.576L223.529 0Z"
          fill="#E74BAC"
        />
      </svg>
    ),
  }),
  Green: createIcon({
    viewBox: "0 0 429 315",
    path: (
      <svg
        width="429"
        height="315"
        viewBox="0 0 429 315"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M215.323 0L0 168.975L215.323 315L429 168.975L215.323 0Z"
          fill="#30D15D"
        />
      </svg>
    ),
  }),
};

// abs. Box with the shapes.
export const BlurShapesBox = (props: BoxProps) => {
  return (
    <Box position="absolute" width="674px" height="479px" {...props}>
      <Box
        position="absolute"
        width="1189px"
        height="100vh"
        left="0px"
        top="0%"
        bottom="0%"
        background="rgba(45, 47, 48, 0.02)"
        backdropFilter="blur(80px)"
        zIndex={1}
      />

      <BlurShapes.Pink
        position="absolute"
        width="372px"
        height="310px"
        left="calc(50% - 372px/2 - 119px)"
        top="186px"
        zIndex={0}
      />

      <BlurShapes.Blue
        position="absolute"
        width="476px"
        height="426px"
        left="calc(50% - 476px/2 + 131px)"
        top="160px"
        zIndex={0}
      />

      <BlurShapes.Green
        position="absolute"
        width="429px"
        height="315px"
        left="calc(50% - 429px/2 + 121.5px)"
        top="324px"
        zIndex={0}
      />
    </Box>
  );
};
