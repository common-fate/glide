import React, { useMemo } from "react";
import { Controller, useFormContext, Validate } from "react-hook-form";
import ReactSelect, { ActionMeta, components, OptionProps } from "react-select";
import { Text } from "@chakra-ui/react";
import {
  useListGroups,
  useGetUsers,
} from "../../../../utils/backend-client/admin/admin";
import { colors } from "../../../../utils/theme/colors";
import { Option } from "../../../../utils/backend-client/types/accesshandler-openapi.yml";
import { ListGroupsSource } from "../../../../utils/backend-client/types";
interface BaseSelectProps {
  fieldName: string;
  rules?: MultiSelectRules;
  isDisabled?: boolean;
  testId?: string;
  onBlurSecondaryAction?: () => void;
}

interface GroupSelectProps extends BaseSelectProps {
  shouldShowGroupMembers?: boolean;
  source?: ListGroupsSource;
}

// UserSelect required defaults to true
export const UserSelect: React.FC<BaseSelectProps> = (props) => {
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
  const { data } = useListGroups({ source: props.source });
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
  validate: Validate<any> | Record<string, Validate<any>> | undefined;
}>;

interface MultiSelectProps extends BaseSelectProps {
  options: Option[];
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
    labelPrefix?: string;
  },
  true
>) => {
  console.log({ children, innerProps });
  return (
    // @ts-ignore
    <div data-testid={innerProps.value}>
      <components.Option {...innerProps}>
        <>
          <Text textStyle={"Body/Medium"}>
            {innerProps?.data.labelPrefix !== undefined && (
              <Text textStyle={"Body/Small"} color="neutrals.500" as={"span"}>
                {innerProps?.data.labelPrefix}
              </Text>
            )}
            {children}
          </Text>
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
          <ReactSelect
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
              multiValue: (provided: any, state: any) => {
                return {
                  ...provided,
                  borderRadius: "20px",
                  background: colors.neutrals[100],
                  // @TODO Hack: I couldn't work out why the layout was overflowing the step container so I added this as a workaround to fix it
                  // this doesn't work on small screens like mobile
                  maxWidth: "400px",
                };
              },
              container: (provided: any, state: any) => {
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
            onChange={(val: any) => {
              // for MultiSelect with 'Select All' option
              // we check if the selected value is 'select all'
              // if true then add all options as value.
              if (shouldAddSelectAllOption) {
                const isAllOptionSelected = val.find(
                  (c: any) => c.label === SELECT_ALL_LABEL
                );

                if (isAllOptionSelected) {
                  onChange(sortedOptions.map((o) => o.value));

                  return;
                }
              }
              onChange(val.map((c: any) => c.value));
            }}
            onBlur={() => {
              rest.onBlurSecondaryAction && rest.onBlurSecondaryAction();
              onBlur();
            }}
            data-testid={rest.testId}
            {...rest}
          />
        );
      }}
    />
  );
};
interface SelectProps extends BaseSelectProps {
  options: {
    value: string;
    label: string;
    description?: string;
  }[];
  id?: string;
}

// single select but uses an array as the underlying data type
export const SelectWithArrayAsValue: React.FC<SelectProps> = ({
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
      name={fieldName}
      render={({ field: { onChange, ref, value, onBlur } }) => {
        const onChangeHandler = (
          option: Option | null,
          actionMeta: ActionMeta<Option>
        ) => {
          onChange(option?.value ? [option?.value] : []);
        };
        return (
          <ReactSelect
            id={id}
            isDisabled={isDisabled}
            //getOptionLabel={(option) => `${option.label}  (${option.value})`}
            options={sortedOptions}
            components={{ Option: CustomOption }}
            onMenuClose={onBlur}
            styles={{
              singleValue: (provided: any, state: any) => {
                return {
                  ...provided,
                  borderRadius: "20px",
                  background: colors.neutrals[100],
                  // @TODO Hack: I couldn't work out why the layout was overflowing the step container so I added this as a workaround to fix it
                  // this doesn't work on small screens like mobile
                  maxWidth: "400px",
                };
              },
              container: (provided: any, state: any) => {
                return {
                  ...provided,
                  // @TODO Hack: I couldn't work out why the layout was overflowing the step container so I added this as a workaround to fix it
                  // this doesn't work on small screens like mobile
                  minWidth: "calc(-20px + 100%)",
                };
              },
              option: (provided: any, state: any) => {
                return {
                  ...provided,
                  background: state.isSelected
                    ? colors.blue[200]
                    : provided.background,
                  color: state.isSelected
                    ? colors.neutrals[800]
                    : provided.color,
                };
              },
            }}
            // ref={ref}
            value={
              value?.length === 1
                ? sortedOptions.find((c) => value[0] === c.value)
                : undefined
            }
            // @ts-ignore
            onChange={onChangeHandler}
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
