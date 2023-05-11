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
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
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

interface ResourceSchemaProperty {
  title: string;
  type: string;
}

interface TargetGroupFilterOperation {
  attribute: string;
  values: string[];
  value: string;
  operationType: string;
}

export const TargetGroupField: React.FC<TargetGroupFieldProps> = (props) => {
  const { targetGroup, fieldSchema } = props;
  const [resources, setResources] = useState<TargetGroupResource[]>([]);
  const { isOpen, onToggle, onClose } = useDisclosure();
  const [isLoading, setIsLoading] = useState<Boolean>(false);

  const { watch, register, control, setValue } = useFormContext<any>();
  const [targetGroupFilterOpts] = watch(["targetgroups"]);

  const [operationType, setOperationType] =
    useState<ResourceFilterOperationTypeEnum>(
      ResourceFilterOperationTypeEnum.IN
    );

  const [filteredResources, setFilterResources] = useState<
    TargetGroupResource[]
  >([]);

  useEffect(() => {
    if (fieldSchema?.resource) {
      adminGetTargetGroupResources(targetGroup.id, fieldSchema.resource).then(
        (data) => {
          setResources(data);
        }
      );
    }
    onClose();
  }, [targetGroup]);

  const fetchFilterResources = async () => {
    const resourceFilter = createResourceFilter(
      targetGroupFilterOpts[targetGroup.id][fieldSchema.id]
    );

    if (!resourceFilter || !fieldSchema.resource) {
      return;
    }

    setIsLoading(true);
    adminFilterTargetGroupResources(
      targetGroup.id,
      fieldSchema.resource,
      resourceFilter
    )
      .then((data) => {
        setFilterResources(data);
        setIsLoading(false);
      })
      .catch(() => setIsLoading(false));
  };

  const createOptions = () => {
    let defaultOptions = [
      { value: "id", label: "Id" },
      { value: "name", label: "Name" },
    ];

    if (fieldSchema.resourceSchema) {
      const properties = fieldSchema.resourceSchema.properties;

      if (properties?.data) {
        let out = Object.entries(
          properties.data as Map<string, ResourceSchemaProperty>
        )
          .filter(([k, v]) => v.type === "string")
          .map(([k, v]) => ({ value: k, label: v.title }));

        defaultOptions.push(...out);
      }
    }

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
        {false ? (
          <Spinner />
        ) : (
          <>
            {!isOpen && (
              <Box>
                <Text
                  textStyle={"Body/Medium"}
                  color="neutrals.500"
                  pl={4}
                  pb={2}
                  pt={2}
                >
                  Available {fieldSchema.title}s
                </Text>
                {resources.length!! ? (
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
                ) : (
                  <Alert status="warning">
                    <AlertIcon />
                    <AlertTitle>No resources synced!</AlertTitle>
                    <AlertDescription>
                      You can manually sync your target resources by running
                      cache sync command
                    </AlertDescription>
                  </Alert>
                )}
              </Box>
            )}
            <Box p={4}>
              <Collapse in={isOpen} animateOpacity>
                <HStack>
                  <Text color={"black"}> Where</Text>
                  <Controller
                    name={`targetgroups.${targetGroup.id}.${fieldSchema.id}.attribute`}
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
                        onChange={(val: any) => {
                          onChange(val?.value);
                        }}
                        value={createOptions().find((c) => c.value === value)}
                      />
                    )}
                  />
                  <Controller
                    name={`targetgroups.${targetGroup.id}.${fieldSchema.id}.operationType`}
                    control={control}
                    defaultValue={"IN"}
                    render={({ field: { value, ref, name, onChange } }) => (
                      <ReactSelect
                        ref={ref}
                        styles={{
                          control: (provided, state) => ({
                            ...provided,
                            width: 120,
                          }),
                        }}
                        value={createFilterOperationOptions().find((c) => {
                          return c.value === value;
                        })}
                        options={createFilterOperationOptions()}
                        onChange={(val) => {
                          onChange(val?.value);
                          setOperationType(
                            val?.value as ResourceFilterOperationTypeEnum
                          );
                        }}
                      />
                    )}
                  />
                  {operationType === ResourceFilterOperationTypeEnum.IN ? (
                    <MultiSelect
                      options={resources.map((r) => {
                        return { label: r.resource.name, value: r.resource.id };
                      })}
                      fieldName={`targetgroups.${targetGroup.id}.${fieldSchema.id}.values`}
                      onBlurSecondaryAction={() => fetchFilterResources()}
                    />
                  ) : (
                    <FormControl>
                      <Input
                        {...register(
                          `targetgroups.${targetGroup.id}.${fieldSchema.id}.value`
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

                      setValue(
                        `targetgroups.${targetGroup.id}.${fieldSchema.id}.attribute`,
                        "id"
                      );

                      setValue(
                        `targetgroups.${targetGroup.id}.${fieldSchema.id}.values`,
                        []
                      );

                      setValue(
                        `targetgroups.${targetGroup.id}.${fieldSchema.id}.operationType`,
                        ResourceFilterOperationTypeEnum.IN
                      );

                      setValue(
                        `targetgroups.${targetGroup.id}.${fieldSchema.id}.value`,
                        ""
                      );

                      setFilterResources([]);
                    }}
                    ml="60%"
                  >
                    Remove Filter
                  </Button>
                )}
              </div>
            </Flex>
          </>
        )}
      </VStack>
    </Box>
  );
};

export const createResourceFilter = (
  operation: TargetGroupFilterOperation
): ResourceFilter => {
  if (!operation.values?.length && operation.value == "") {
    return [];
  }

  return [
    {
      operationType: operation.operationType as ResourceFilterOperationTypeEnum,
      attribute: operation.attribute,
      ...(operation.operationType === ResourceFilterOperationTypeEnum.IN
        ? {
            values: operation.values || [],
          }
        : {
            value: operation.value,
          }),
    },
  ];
};
