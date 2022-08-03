import React, { useMemo } from "react";
import { useFormContext, Controller } from "react-hook-form";
import Select from "react-select";
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

export const GroupSelect: React.FC<SelectProps> = (props) => {
  const { data } = useGetGroups();
  const options = useMemo(() => {
    return (
      data?.groups.map((g) => {
        return { value: g.id, label: g.name };
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
            onBlur={() => trigger(fieldName)}
            data-testid={rest.testId}
            {...rest}
          />
        );
      }}
    />
  );
};
