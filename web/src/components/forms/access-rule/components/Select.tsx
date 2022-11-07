import React, { useMemo } from "react";
import { Controller, useFormContext } from "react-hook-form";
import Select, { components, OptionProps } from "react-select";
import { Text } from "@chakra-ui/react";
import {
  useListGroups,
  useGetUsers,
} from "../../../../utils/backend-client/admin/admin";
import { colors } from "../../../../utils/theme/colors";
interface SelectProps {
  fieldName: string;
  rules?: MultiSelectRules;
  isDisabled?: boolean;
  testId?: string;
  onBlurSecondaryAction?: () => void;
}

interface GroupSelectProps extends SelectProps {
  shouldShowGroupMembers?: boolean;
}

// UserSelect required defaults to true
export const UserSelect: React.FC<SelectProps> = (props) => {
  const { data } = useGetUsers();
  const options = useMemo(() => {
    return (
      data?.users
        .map((u) => {
          return { value: u.id, label: u.email };
        })
        .sort((a, b) => a.label.localeCompare(b.label)) ?? []
    );
  }, [data]);
  return <MultiSelect id="user-select" options={options} {...props} />;
};

export const GroupSelect: React.FC<GroupSelectProps> = (props) => {
  const { shouldShowGroupMembers = false } = props;
  const { data } = useListGroups();
  const options = useMemo(() => {
    return (
      data?.groups
        .map((g) => {
          const totalMembersInGroup =
            g.memberCount <= 1
              ? `${g.memberCount} member`
              : `${g.memberCount} members`;

          return {
            value: g.id,
            label: shouldShowGroupMembers
              ? `${g.name} (${totalMembersInGroup})`
              : g.name,
          };
        })
        .sort((a, b) => a.label.localeCompare(b.label)) ?? []
    );
  }, [data, shouldShowGroupMembers]);
  return <MultiSelect id={props.testId} options={options} {...props} />;
};

type MultiSelectRules = Partial<{
  required: boolean;
  minLength: number;
}>;
interface MultiSelectProps extends SelectProps {
  options: {
    value: string;
    label: string;
    description?: string;
  }[];
  id?: string;
  shouldAddSelectAllOption?: boolean;
}

export const CustomOption = ({
  children,
  ...innerProps
}: OptionProps<
  {
    value: string;
    label: string;
    description?: string;
  },
  true
>) => {
  return (
    // @ts-ignore
    <div data-testid={innerProps.value}>
      <components.Option {...innerProps}>
        <>
          {children}
          {innerProps?.data.description && (
            <Text>{innerProps.data.description}</Text>
          )}
          {
            // @ts-ignore
            <Text>{innerProps.value}</Text>
          }
        </>
      </components.Option>
    </div>
  );
};

const SELECT_ALL_LABEL = "Select all";
const SELECT_ALL_OPTION = { label: SELECT_ALL_LABEL, value: "" };

export const MultiSelect: React.FC<MultiSelectProps> = ({
  options,
  fieldName,
  rules,
  isDisabled,
  id,
  shouldAddSelectAllOption = false,
  ...rest
}) => {
  const { control, trigger } = useFormContext();
  // sort the options alphabetically
  const sortedOptions = useMemo(
    () => options.sort((a, b) => a.label.localeCompare(b.label)),
    [options]
  );

  return (
    <Controller
      control={control}
      rules={{ ...rules }}
      defaultValue={[]}
      name={fieldName}
      render={({ field: { onChange, ref, value, onBlur } }) => {
        return (
          <Select
            id={id}
            isDisabled={isDisabled}
            //getOptionLabel={(option) => `${option.label}  (${option.value})`}
            options={[
              ...(shouldAddSelectAllOption
                ? [SELECT_ALL_OPTION, ...sortedOptions]
                : sortedOptions),
            ]}
            components={{ Option: CustomOption }}
            isMulti
            onMenuClose={onBlur}
            styles={{
              multiValue: (provided, state) => {
                return {
                  ...provided,
                  borderRadius: "20px",
                  background: colors.neutrals[100],
                  // @TODO Hack: I couldn't work out why the layout was overflowing the step container so I added this as a workaround to fix it
                  // this doesn't work on small screens like mobile
                  maxWidth: "400px",
                };
              },
              container: (provided, state) => {
                return {
                  ...provided,
                  // @TODO Hack: I couldn't work out why the layout was overflowing the step container so I added this as a workaround to fix it
                  // this doesn't work on small screens like mobile
                  minWidth: "calc(-20px + 100%)",
                };
              },
            }}
            // ref={ref}
            value={sortedOptions.filter((c) => value.includes(c.value))}
            onChange={(val) => {
              // for MultiSelect with 'Select All' option
              // we check if the selected value is 'select all'
              // if true then add all options as value.
              if (shouldAddSelectAllOption) {
                const isAllOptionSelected = val.find(
                  (c) => c.label === SELECT_ALL_LABEL
                );

                if (isAllOptionSelected) {
                  onChange(sortedOptions.map((o) => o.value));

                  return;
                }
              }
              onChange(val.map((c) => c.value));
            }}
            onBlur={() => {
              void trigger(fieldName);
              rest.onBlurSecondaryAction && rest.onBlurSecondaryAction();
            }}
            data-testid={rest.testId}
            {...rest}
          />
        );
      }}
    />
  );
};
// MultiSelectOptions is the same as MultiSelect except that it sets the field value to the array of options rather than an array of values
export const MultiSelectOptions: React.FC<MultiSelectProps> = ({
  options,
  fieldName,
  rules,
  isDisabled,
  id,
  ...rest
}) => {
  const { control, trigger } = useFormContext();
  // sort the options alphabetically
  const sortedOptions = useMemo(
    () => options.sort((a, b) => a.label.localeCompare(b.label)),
    [options]
  );
  return (
    <Controller
      control={control}
      rules={{ ...rules }}
      defaultValue={[]}
      name={fieldName}
      render={({ field: { onChange, ref, value } }) => {
        return (
          <Select
            id={id}
            isDisabled={isDisabled}
            options={sortedOptions}
            components={{ Option: CustomOption }}
            isMulti
            styles={{
              multiValue: (provided, state) => {
                return {
                  ...provided,
                  borderRadius: "20px",
                  background: colors.neutrals[100],
                };
              },
              container: (provided, state) => {
                return {
                  ...provided,
                  minWidth: "100%",
                };
              },
            }}
            ref={ref}
            value={value ?? []}
            onChange={(val) => {
              onChange(val);
            }}
            onBlur={() => {
              void trigger(fieldName);
              rest.onBlurSecondaryAction && rest.onBlurSecondaryAction();
            }}
            data-testid={rest.testId}
            {...rest}
          />
        );
      }}
    />
  );
};
