import {
  Box,
  Button,
  Collapse,
  Flex,
  FormLabel,
  HStack,
  VStack,
  Input,
  Spinner,
  Text,
  useDisclosure,
  Wrap,
  WrapItem,
  Container,
} from "@chakra-ui/react";
import React, { useEffect, useState } from "react";

import { useAdminGetTargetGroupResources } from "../../../../utils/backend-client/admin/admin";
import ReactSelect from "react-select";
import {
  ResourceFilter,
  ResourceFilterOperationTypeEnum,
  TargetGroupResource,
  TargetGroup,
} from "../../../../utils/backend-client/types";
import { adminFilterTargetGroupResources } from "../../../../utils/backend-client/default/default";
import { DynamicOption } from "../../../../components/DynamicOption";
import { CloseIcon, StarIcon } from "@chakra-ui/icons";
import { FilterIcon } from "../../../../components/icons/Icons";

interface TargetGroupFieldProps {
  targetGroup?: TargetGroup;
  resourceType: string;
}

export const TargetGroupField: React.FC<TargetGroupFieldProps> = (props) => {
  const { resourceType = "Account" } = props;
  const { data } = useAdminGetTargetGroupResources(
    "common-fate-aws",
    "account"
  );

  const [selectedFilters, setSelectedFilters] = useState<Number>(0);
  const { isOpen, onToggle } = useDisclosure();

  const [selectedAttribute, setSelectedAttribute] = useState<string>("");
  const [selectedOperationType, setSelectedOperationType] = useState<
    ResourceFilterOperationTypeEnum | ""
  >("");

  const [filteredResources, setFilteredResources] = useState<
    TargetGroupResource[]
  >([]);

  useEffect(() => {
    const resourceFilter = createResourceFilter();
    if (!resourceFilter) {
      return;
    }

    console.log("the resourceFilter is", resourceFilter);

    adminFilterTargetGroupResources("abcd", "reourceType", resourceFilter).then(
      (resp) => {
        console.log("thee resp is", resp);
        setFilteredResources(resp);
      }
    );
  }, [selectedOperationType]);

  const createResourceFilter = (): ResourceFilter | void => {
    if (selectedOperationType == "") {
      return;
    }

    return [
      {
        operationType: selectedOperationType,
        attribute: selectedAttribute,
        ...(selectedOperationType === ResourceFilterOperationTypeEnum.IN
          ? { values: ["abc", "efg"] }
          : { value: "mockaccount" }),
      },
    ];
  };

  if (!data) {
    return <Spinner />;
  }

  const createOptions = () => {
    const options = [
      { value: "id", label: "id" },
      { value: "name", label: "name" },
    ];

    return options;
  };

  return (
    <>
      <VStack
        w="100%"
        align={"flex-start"}
        spacing={4}
        // key={k}
        p={4}
        rounded="md"
        border="1px solid"
        borderColor="gray.300"
      >
        <Box>
          <Text textStyle={"Body/Medium"} color="neutrals.500">
            Available Accounts
          </Text>
          <Wrap>
            {data.slice(0, 10).map((opt) => {
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

        <Container>
          <Collapse in={isOpen} animateOpacity>
            <HStack>
              <Text> Where</Text>
              <ReactSelect
                styles={{
                  control: (provided, state) => ({
                    ...provided,
                    width: 120,
                  }),
                }}
                defaultValue={{ value: "id", label: "id" }}
                options={createOptions()}
                onChange={(selectedOption) => {
                  if (selectedOption?.value) {
                    setSelectedAttribute(selectedOption.value);
                  }
                }}
              />
              <ReactSelect
                styles={{
                  control: (provided, state) => ({
                    ...provided,
                    width: 120,
                  }),
                }}
                defaultValue={{
                  value: ResourceFilterOperationTypeEnum.IN as string,
                  label: "in",
                }}
                options={Object.entries(ResourceFilterOperationTypeEnum).map(
                  ([value, label]) => ({
                    value,
                    label: label.toLowerCase(),
                  })
                )}
                onChange={(selectedOption) => {
                  if (selectedOption?.value) {
                    setSelectedOperationType(
                      selectedOption.value as ResourceFilterOperationTypeEnum
                    );
                  }
                }}
              />
              {selectedOperationType === ResourceFilterOperationTypeEnum.IN ? (
                <ReactSelect
                  styles={{
                    control: (provided, state) => ({
                      ...provided,
                      width: 200,
                    }),
                  }}
                  options={[
                    { label: "abc", value: "abc" },
                    { label: "appl", value: "appl" },
                  ]}
                />
              ) : (
                <Input />
              )}
            </HStack>
            <Box overflow="auto" h="200px">
              <Text>Preview</Text>
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
            <Flex justifyContent={"space-between"}>
              <Text>
                {selectedFilters == 0
                  ? "No Filters Applied"
                  : `${+selectedFilters} Filters Applied`}
              </Text>
              {!isOpen ? (
                <Button
                  leftIcon={<FilterIcon />}
                  size="m"
                  variant="ghost"
                  onClick={onToggle}
                  ml="60%"
                >
                  Apply Filter
                </Button>
              ) : (
                <Button
                  leftIcon={<CloseIcon boxSize={2} />}
                  size="m"
                  color="red.400"
                  variant="ghost"
                  onClick={onToggle}
                  ml="60%"
                >
                  Remove Filter
                </Button>
              )}
            </Flex>
          </Box>
        </Container>
      </VStack>
    </>
  );
};
