import { EditIcon } from "@chakra-ui/icons";
import {
  Flex,
  HStack,
  Menu,
  MenuIcon,
  MenuItem,
  MenuList,
  SkeletonText,
  Text,
  Tooltip,
} from "@chakra-ui/react";
import { useMemo } from "react";
import { Link, MakeGenerics, useNavigate, useSearch } from "react-location";
import { Column } from "react-table";
import { makeLookupAccessRuleRequestLink } from "../../pages/access";
import { useUserGetAccessRuleApprovers } from "../../utils/backend-client/end-user/end-user";
import {
  AccessRule,
  AccessRuleStatus,
  KeyValue,
  LookupAccessRule,
} from "../../utils/backend-client/types";
import { durationString } from "../../utils/durationString";
import { ProviderIcon } from "../icons/providerIcon";
import { TableRenderer } from "./TableRenderer";

// Note: I made this type because the table column types don't seem to work with complex data
// We need access to the selectabeloptions when redirecting to the request page
interface ExtendedAR extends AccessRule {
  selectableWithOptionValues?: KeyValue[];
}
export const SelectRuleTable = ({ rules }: { rules: LookupAccessRule[] }) => {
  const navigate = useNavigate();

  const cols = useMemo<Column<ExtendedAR>[]>(
    () => [
      {
        accessor: "name",
        Header: "Name",
        Cell: ({ cell }) => {
          return (
            <Text textStyle="Body/Small" color="neutrals.700" as="a">
              {cell.value}
            </Text>
          );
        },
      },
      {
        accessor: "description",
        Header: "Description",

        Cell: ({ cell }) => {
          return (
            // Truncates the text if it is long, full description is in the tooltip
            <Tooltip label={cell.value} aria-label="description">
              <Text
                textStyle="Body/Small"
                color="neutrals.700"
                noOfLines={1}
                maxWidth="200px"
              >
                {cell.value}
              </Text>
            </Tooltip>
          );
        },
      },
      // {
      //   accessor: "status",
      //   Header: "Status",
      //   Cell: ({ cell }) => {
      //     return (
      //       <StatusCell
      //         value={cell.value}
      //         danger={AccessRuleStatus.ARCHIVED}
      //         success={AccessRuleStatus.ACTIVE}
      //         textStyle="Body/Small"
      //       />
      //     );
      //   },
      // },
      {
        accessor: "target",
        Header: "Details",
        Cell: ({ cell }) => {
          return (
            <HStack>
              <ProviderIcon shortType={cell.value.provider.type} />

              <Text
                color="neutrals.700"
                textStyle="Body/Small"
                whiteSpace={"nowrap"}
              >
                {durationString(
                  cell.row.original.timeConstraints.maxDurationSeconds
                )}
              </Text>
            </HStack>
          );
        },
      },
      {
        accessor: "id",
        id: "approvers",
        Header: "Approval Required?",
        Cell: ({ cell }) => {
          const { data, isValidating } = useUserGetAccessRuleApprovers(
            cell.row.original.id
          );

          return (
            <HStack>
              <Text
                color="neutrals.700"
                textStyle="Body/Small"
                whiteSpace={"nowrap"}
              >
                {isValidating || !data ? (
                  <SkeletonText w="12ch" noOfLines={1} h="12px" />
                ) : data?.users?.length > 0 ? (
                  "Yes"
                ) : (
                  "No"
                )}
              </Text>
            </HStack>
          );
        },
      },

      {
        accessor: "id",
        Header: "",
        id: "actions",
        Cell: ({ cell }) => {
          return (
            <Menu>
              <MenuList>
                <Link to={"/admin/access-rules/" + cell.value}>
                  <MenuItem as="a">
                    <MenuIcon mr={2} color="neutrals.500">
                      <EditIcon />
                    </MenuIcon>
                    Edit Rule
                  </MenuItem>
                </Link>
              </MenuList>
            </Menu>
          );
        },
      },
    ],
    []
  );

  return (
    <Flex mt={8}>
      {TableRenderer<ExtendedAR>({
        columns: cols,
        data: rules.map((d) => {
          return {
            ...d.accessRule,
            selectableWithOptionValues: d.selectableWithOptionValues,
          };
        }),
        emptyText: "No access rules",
        rowProps: (rule) => ({
          cursor: "pointer",
          onClick: (e) => {
            const { selectableWithOptionValues, ...accessRule } = rule.original;
            e.preventDefault();
            navigate(
              makeLookupAccessRuleRequestLink({
                accessRule,
                selectableWithOptionValues,
              })
            );
          },
        }),
      })}
    </Flex>
  );
};
