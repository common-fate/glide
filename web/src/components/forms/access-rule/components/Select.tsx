import React, { useEffect, useMemo, useState } from "react";
import { Controller, useFormContext, Validate } from "react-hook-form";
import ReactSelect, { ActionMeta, components, OptionProps } from "react-select";
import { Box, Text } from "@chakra-ui/react";
import {
  adminListGroups,
  adminListUsers,
} from "../../../../utils/backend-client/admin/admin";
import { colors } from "../../../../utils/theme/colors";

import {
  AdminListGroupsSource,
  Group,
  User,
} from "../../../../utils/backend-client/types";

interface Option {
  label: string;
  value: string;
}
interface BaseSelectProps {
  fieldName: string;
  rules?: MultiSelectRules;
  isDisabled?: boolean;
  testId?: string;
  onBlurSecondaryAction?: () => void;
}

interface GroupSelectProps extends BaseSelectProps {
  shouldShowGroupMembers?: boolean;
  source?: AdminListGroupsSource;
  onBlurSecondaryAction?: () => void;
}

// UserSelect required defaults to true
export const UserSelect: React.FC<BaseSelectProps> = (props) => {
  const [items, setItems] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  async function fetchUsers(nextInput?: string) {
    // if there are items, and no nextInput then we are done
    if (!nextInput && items.length !== 0) {
      setIsLoading(false);
      return;
    }

    const { users, next } = await adminListUsers({
      // this is to accomodate for type discrepancy on first call
      nextToken: nextInput ? nextInput : undefined,
    });

    if (users) {
      users && setItems((currItems) => [...currItems, ...users]);
      if (next) {
        void fetchUsers(next);
      } else {
        setIsLoading(false);
      }
    }
  }

  useEffect(() => {
    void fetchUsers();
    return () => {
      setItems([]);
      setIsLoading(true);
    };
  }, []);

  const options = useMemo(
    () =>
      items
        .map((u) => {
          return { value: u.id, label: u.email };
        })
        // filter out dupes:
        .filter((v, i, a) => a.findIndex((t) => t.value === v.value) === i)
        .sort((a, b) => a.label.localeCompare(b.label)) ?? [],
    [items]
  );
  return (
    <MultiSelect
      id="user-select"
      isLoading={isLoading}
      options={options}
      {...props}
    />
  );
};

export const GroupSelect: React.FC<GroupSelectProps> = (props) => {
  const { shouldShowGroupMembers = false } = props;

  const [items, setItems] = useState<Group[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  // const fetchGroups = async (nextInput?: string) => { // recreate but mark as void
  async function fetchGroups(nextInput?: string) {
    // if there are items, and no nextInput then we are done
    if (!nextInput && items.length !== 0) {
      setIsLoading(false);
      return;
    }

    const { groups, next } = await adminListGroups({
      // this is to accomodate for type discrepancy on first call
      nextToken: nextInput ? nextInput : undefined,
      source: props.source,
    });
    if (groups) {
      groups && setItems((currItems) => [...currItems, ...groups]);
      if (next) {
        void fetchGroups(next);
      } else {
        setIsLoading(false);
      }
    }
  }

  useEffect(() => {
    void fetchGroups();
    return () => {
      setItems([]);
      setIsLoading(true);
    };
  }, []);

  const options = useMemo(() => {
    return (
      items
        .map((g) => {
          const totalMembersInGroup =
            g.memberCount <= 1
              ? g.memberCount == 0
                ? "No members"
                : `${g.memberCount} member`
              : `${g.memberCount} members`;
          return {
            value: g.id,
            label: shouldShowGroupMembers
              ? `${g.name} (${totalMembersInGroup})`
              : g.name,
          };
        })
        // let's run a quick filter to remove any dupes
        .filter((v, i, a) => a.findIndex((t) => t.value === v.value) === i)
        .sort((a, b) => a.label.localeCompare(b.label)) ?? []
    );
  }, [items, shouldShowGroupMembers]);
  return (
    <MultiSelect
      id={props.testId}
      options={options}
      {...props}
      onBlurSecondaryAction={props.onBlurSecondaryAction}
      isLoading={isLoading}
    />
  );
};

type MultiSelectRules = Partial<{
  required: boolean;
  minLength: number;
  // @ts-ignore; the CI has randomly started throwing an error here despite typings being fine in IDE, CI may be using an incorrect version of react-hook-form
  validate: Validate<any> | Record<string, Validate<any>> | undefined;
}>;

interface MultiSelectProps extends BaseSelectProps {
  options: Option[];
  id?: string;
  shouldAddSelectAllOption?: boolean;
  isLoading?: boolean;
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
  console.debug({ children, innerProps });
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
  isLoading,
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
            key={id}
            isLoading={isLoading}
            inputId={id + "-input"}
            isDisabled={isDisabled}
            //getOptionLabel={(option) => `${option.label}  (${option.value})`}
            options={[
              ...(shouldAddSelectAllOption
                ? [SELECT_ALL_OPTION, ...sortedOptions]
                : sortedOptions),
            ]}
            components={{
              Option: CustomOption,
              ClearIndicator: (props) => (
                <Box data-testid="select-input-deselect-all">
                  <components.ClearIndicator {...props} />
                </Box>
              ),
            }}
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
                  // minWidth: "calc(-20px + 100%)",
                  position: "relative",
                  width: "100%",
                  display: "flex",
                  flexGrow: 1,
                };
              },
              control: (provided: any, state: any) => {
                return {
                  ...provided,

                  borderColor: "#E5E5E5", // neutrals.300
                  width: "100%",
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
            data-testid="select-input"
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
              control: (provided: any, state: any) => {
                return {
                  ...provided,
                  borderColor: "#E5E5E5", // neutrals.300
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
