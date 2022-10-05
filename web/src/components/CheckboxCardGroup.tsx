import {
  Box,
  BoxProps,
  Checkbox,
  Stack,
  StackProps,
  useCheckbox,
  useCheckboxGroup,
  UseCheckboxGroupProps,
  UseCheckboxProps,
  useId,
  useStyleConfig,
} from "@chakra-ui/react";
import React from "react";

type CheckboxCardGroupProps = StackProps & UseCheckboxGroupProps;

export const CheckboxCardGroup = (props: CheckboxCardGroupProps) => {
  const { children, defaultValue, value, onChange, ...rest } = props;
  const { getCheckboxProps } = useCheckboxGroup({
    defaultValue,
    value,
    onChange,
  });

  const cards = React.useMemo(
    () =>
      React.Children.toArray(children)
        .filter<React.ReactElement<RadioCardProps>>(React.isValidElement)
        .map((card) => {
          return React.cloneElement(card, {
            checkboxProps: getCheckboxProps({
              value: card.props.value,
            }),
          });
        }),
    [children, getCheckboxProps]
  );

  return <Stack {...rest}>{cards}</Stack>;
};

interface RadioCardProps extends BoxProps {
  value: string;
  checkboxProps?: UseCheckboxProps;
}

export const CheckboxCard = (props: RadioCardProps) => {
  const { checkboxProps, children, ...rest } = props;
  const { getInputProps, getCheckboxProps, getLabelProps, state } =
    useCheckbox(checkboxProps);
  const id = useId(undefined, "checkbox-card");
  const styles = useStyleConfig("RadioCard", props);

  return (
    <Box
      as="label"
      cursor="pointer"
      {...getLabelProps()}
      sx={{
        ".focus-visible + [data-focus]": {
          boxShadow: "outline",
          zIndex: 1,
        },
      }}
    >
      <input {...getInputProps()} aria-labelledby={id} />
      <Box
        {...getCheckboxProps()}
        {...rest}
        sx={{
          // ...styles,
          borderWidth: "1px",
          borderRadius: "lg",
          p: "4",
          bg: "bg-surface",
          transitionProperty: "common",
          transitionDuration: "normal",
          _hover: { borderColor: "neutrals.200" },
          _checked: {
            // borderColor: 'brandBlue.500',
            borderWidth: "1px solid #2e7fff",
            // border: '1px solid #2e7fff',
            boxShadow: `0px 0px 0px 1px #2e7fff`,
            // boxShadow:
            // `0px 0px 0px 1px ${transparentize(`brand.200`, 1.0)(props.theme)}`
          },
          _dark: {
            borderColor: "neutrals.300",
            _checked: {
              borderColor: "brandBlue.200",
            },
          },
        }}
      >
        <Stack direction="row">
          <Box flex="1">{children}</Box>
          <Checkbox
            pointerEvents="none"
            isFocusable={false}
            isChecked={state.isChecked}
            alignSelf="start"
          />
        </Stack>
      </Box>
    </Box>
  );
};
