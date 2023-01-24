import {
  mode,
  StyleFunctionProps,
  transparentize,
} from "@chakra-ui/theme-tools";

const baseStyle = (props: StyleFunctionProps) => ({
  borderWidth: "1px",
  borderRadius: "lg",
  p: "4",
  bg: "bg-surface",
  transitionProperty: "common",
  transitionDuration: "normal",
  _hover: { borderColor: mode("neutrals.300", "neutrals.600")(props) },
  _checked: {
    borderColor: mode("brand.500", "brand.200")(props),
    boxShadow: mode(
      `0px 0px 0px 1px ${transparentize(`brand.500`, 1.0)(props.theme)}`,
      `0px 0px 0px 1px ${transparentize(`brand.200`, 1.0)(props.theme)}`
    )(props),
  },
});

export default {
  baseStyle,
};
