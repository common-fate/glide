import React, { useMemo } from "react";
import { useFormContext, Controller } from "react-hook-form";
import Select, { components, OptionProps } from "react-select";
import {
  useGetUsers,
  useGetGroups,
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
      data?.users.map((u) => {
        return { value: u.id, label: u.email };
      }) ?? []
    );
  }, [data]);
  return <MultiSelect id="user-select" options={options} {...props} />;
};

export const GroupSelect: React.FC<GroupSelectProps> = (props) => {
  const { shouldShowGroupMembers = false } = props;
  const { data } = useGetGroups();
  const options = useMemo(() => {
    return (
      data?.groups.map((g) => {
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
      }) ?? []
    );
  }, [data]);
  return <MultiSelect id={props.testId} options={options} {...props} />;
};
type MultiSelectRules = Partial<{
  required: boolean;
}>;
interface MultiSelectProps extends SelectProps {
  options: {
    value: string;
    label: string;
  }[];
  id?: string;
}

export const CustomOption = ({
  children,
  ...innerProps
}: OptionProps<
  {
    value: string;
    label: string;
  },
  true
>) => (
  <div data-testid={innerProps.value}>
    <components.Option {...innerProps}>{children}</components.Option>
  </div>
);

const MultiSelect: React.FC<MultiSelectProps> = ({
  options,
  fieldName,
  rules,
  isDisabled,
  id,
  ...rest
}) => {
  const { control, trigger } = useFormContext();

  return (
    <Controller
      control={control}
      rules={{ required: true, minLength: 1, ...rules }}
      defaultValue={[]}
      name={fieldName}
      render={({ field: { onChange, ref, value } }) => {
        return (
          <Select
            id={id}
            isDisabled={isDisabled}
            options={options}
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
            value={options.filter((c) => value.includes(c.value))}
            onChange={(val) => onChange(val.map((c) => c.value))}
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
