import {
  Box,
  Button,
  Collapse,
  Flex,
  HStack,
  VStack,
  Input,
  Spinner,
  Text,
  useDisclosure,
  Wrap,
  StackDivider,
  FormControl,
} from "@chakra-ui/react";
import React, { useEffect, useState } from "react";

import { adminGetTargetGroupResources } from "../../../../utils/backend-client/admin/admin";
import ReactSelect from "react-select";
import {
  ResourceFilter,
  ResourceFilterOperationTypeEnum,
  TargetGroupResource,
  TargetGroup,
  TargetGroupSchemaArgument,
} from "../../../../utils/backend-client/types";
import { adminFilterTargetGroupResources } from "../../../../utils/backend-client/default/default";
import { DynamicOption } from "../../../../components/DynamicOption";
import { CloseIcon } from "@chakra-ui/icons";
import { FilterIcon } from "../../../../components/icons/Icons";
import { MultiSelect } from "./Select";
import { useFormContext, Controller } from "react-hook-form";

interface TargetGroupFieldProps {
  targetGroup: TargetGroup;
  fieldSchema: TargetGroupSchemaArgument;
}

interface SchemaProperty {
  title: string;
  type: string;
  description: string;
}

export const TargetGroupField: React.FC<TargetGroupFieldProps> = (props) => {
  const { targetGroup, fieldSchema } = props;
  const [resources, setResources] = useState<TargetGroupResource[]>([]);
  const { isOpen, onToggle } = useDisclosure();
  const [isLoading, setIsLoading] = useState<Boolean>(false);
  const [filteredResources, setFilteredResources] = useState<
    TargetGroupResource[]
  >([]);

  const { watch, register, control, setValue } = useFormContext<any>();
  const [targetGroupFilterOpts] = watch(["targetgroup"]);
  const [operationType, setOperationType] =
    useState<ResourceFilterOperationTypeEnum>(
      ResourceFilterOperationTypeEnum.IN
    );

  useEffect(() => {
    if (fieldSchema?.resource) {
      adminGetTargetGroupResources(targetGroup.id, fieldSchema.resource).then(
        (data) => setResources(data)
      );
    }
  }, []);

  const fetchFilterResources = async () => {
    const resourceFilter = createResourceFilter();

    console.log("the resrouceFiler is", resourceFilter);

    if (!resourceFilter) {
      return;
    }

    if (!fieldSchema.resource) {
      return;
    }

    setIsLoading(true);
    const resp = await adminFilterTargetGroupResources(
      targetGroup.id,
      fieldSchema.resource,
      resourceFilter
    );

    setIsLoading(false);
    setFilteredResources(resp);
  };

  const createResourceFilter = (): ResourceFilter | void => {
    return [
      {
        operationType: operationType,
        attribute:
          targetGroupFilterOpts[targetGroup.id][fieldSchema.id].attribute,
        ...(operationType === ResourceFilterOperationTypeEnum.IN
          ? {
              values:
                targetGroupFilterOpts[targetGroup.id][fieldSchema.id]
                  .withSelectable || [],
            }
          : {
              value:
                targetGroupFilterOpts[targetGroup.id][fieldSchema.id].withValue,
            }),
      },
    ];
  };

  // if (!resources) {
  //   return <Spinner />;
  // }

  // TODO: Need to instead use the values returned by TargetField Schema.
  const createOptions = () => {
    const defaultOptions = [
      { value: "id", label: "id" },
      { value: "name", label: "name" },
    ];

    return defaultOptions;
  };

  const createFilterOperationOptions = () => {
    return Object.entries(ResourceFilterOperationTypeEnum).map(
      ([value, label]) => ({
        value,
        label: label.toLowerCase(),
      })
    );
  };

  return (
    <Box>
      <Box pt={4}>
        <Text color={"black"}>{fieldSchema.title} </Text>
        <Text size={"sm"}>{fieldSchema.description} </Text>
      </Box>
      <VStack
        divider={<StackDivider borderColor="gray.200" />}
        w="600px"
        align={"flex-start"}
        spacing={4}
        rounded="md"
        border="1px solid"
        bgColor={"white"}
        borderColor="gray.300"
      >
        {!isOpen && (
          <Box>
            <Text
              textStyle={"Body/Medium"}
              color="neutrals.500"
              pl={4}
              pb={2}
              pt={2}
            >
              Available Accounts
            </Text>
            <Box height={"200px"} overflow={"hidden"}>
              <Wrap>
                {resources.map((opt) => {
                  return (
                    <DynamicOption
                      key={"cp-" + opt.resource.id}
                      label={opt.resource.name}
                      value={opt.resource.id}
                    />
                  );
                })}
              </Wrap>
            </Box>
          </Box>
        )}
        <Box p={4}>
          <Collapse in={isOpen} animateOpacity>
            <HStack>
              <Text color={"black"}> Where</Text>
              <Controller
                name={`targetgroup.${targetGroup.id}.${fieldSchema.id}.attribute`}
                control={control}
                defaultValue={"id"}
                render={({ field: { value, ref, name, onChange } }) => (
                  <ReactSelect
                    ref={ref}
                    styles={{
                      control: (provided, state) => ({
                        ...provided,
                        width: 120,
                      }),
                    }}
                    options={createOptions()}
                    onChange={(val: any) => onChange(val?.value)}
                    value={createOptions().find((c) => c.value === value)}
                  />
                )}
              />
              <ReactSelect
                styles={{
                  control: (provided, state) => ({
                    ...provided,
                    width: 120,
                  }),
                }}
                options={createFilterOperationOptions()}
                onChange={(e) =>
                  setOperationType(e?.value as ResourceFilterOperationTypeEnum)
                }
              />
              {operationType === ResourceFilterOperationTypeEnum.IN ? (
                <MultiSelect
                  options={resources.map((r) => {
                    return { label: r.resource.name, value: r.resource.id };
                  })}
                  fieldName={`targetgroup.${targetGroup.id}.${fieldSchema.id}.withSelectable`}
                  onBlurSecondaryAction={() => fetchFilterResources()}
                />
              ) : (
                <FormControl>
                  <Input
                    {...register(
                      `targetgroup.${targetGroup.id}.${fieldSchema.id}.withValue`
                    )}
                    onBlur={() => fetchFilterResources()}
                  />
                </FormControl>
              )}
            </HStack>
            <Box overflow="auto" h="200px" pt={4}>
              <Text>Preview</Text>
              {isLoading && <Spinner />}
              {filteredResources.map((r) => {
                return (
                  <DynamicOption
                    key={"cp-preview-" + r.resource.id}
                    label={r.resource.name}
                    value={r.resource.id}
                  />
                );
              })}
            </Box>
          </Collapse>
          <Box>
            <Flex justifyContent={"space-between"}></Flex>
          </Box>
        </Box>
        <Flex
          justifyContent={"space-between"}
          alignItems={"flex-start"}
          pt={"4px"}
          pb={"4px"}
        >
          <Text>{0} Filters Applied</Text>
          <div>
            {!isOpen ? (
              <Button
                leftIcon={<FilterIcon />}
                size="m"
                variant="ghost"
                onClick={onToggle}
              >
                Apply Filter
              </Button>
            ) : (
              <Button
                leftIcon={<CloseIcon boxSize={2} />}
                size="m"
                color="red.400"
                variant="ghost"
                onClick={() => {
                  onToggle();
                  setValue("operationType", ResourceFilterOperationTypeEnum.IN);
                  setValue("withSelectable", []);
                  setValue("withValue", "");
                  setValue("attribute", "id");

                  setFilteredResources([]);
                }}
                ml="60%"
              >
                Remove Filter
              </Button>
            )}
          </div>
        </Flex>
      </VStack>
    </Box>
  );
};
