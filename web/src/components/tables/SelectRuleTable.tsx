import { ChevronDownIcon, EditIcon, SmallAddIcon } from "@chakra-ui/icons";
import {
  Button,
  Flex,
  HStack,
  Menu,
  MenuButton,
  MenuIcon,
  MenuItem,
  MenuItemOption,
  MenuList,
  MenuOptionGroup,
  SkeletonText,
  Text,
  Tooltip,
} from "@chakra-ui/react";
import { format } from "date-fns";
import { useMemo, useState } from "react";
import { Link, MakeGenerics, useNavigate, useSearch } from "react-location";
import { Column } from "react-table";
import { usePaginatorApi } from "../../utils/usePaginatorApi";
import { useAdminListAccessRules } from "../../utils/backend-client/admin/admin";
import {
  AccessRule,
  AccessRuleDetail,
  AccessRuleStatus,
} from "../../utils/backend-client/types";
import { durationString } from "../../utils/durationString";
import { ProviderIcon } from "../icons/providerIcon";
import RuleConfigModal from "../modals/RuleConfigModal";
import { StatusCell } from "../StatusCell";
import { UserAvatarDetails } from "../UserAvatar";
import { TableRenderer } from "./TableRenderer";
import { useUserGetAccessRuleApprovers } from "../../utils/backend-client/end-user/end-user";

type MyLocationGenerics = MakeGenerics<{
  Search: {
    status?: Lowercase<AccessRuleStatus>;
  };
}>;

export const SelectRuleTable = ({ rules }: { rules: AccessRule[] }) => {
  const search = useSearch<MyLocationGenerics>();
  const navigate = useNavigate<MyLocationGenerics>();

  const status = useMemo(
    () =>
      search.status !== undefined
        ? (search.status.toUpperCase() as AccessRuleStatus)
        : AccessRuleStatus.ACTIVE,
    [search]
  );

  const [selectedAccessRule, setSelectedAccessRule] = useState<AccessRule>();

  const cols = useMemo<Column<AccessRule>[]>(
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
      {TableRenderer<AccessRule>({
        columns: cols,
        data: rules,
        emptyText: "No access rules",
        rowProps: (rule) => ({
          cursor: "pointer",
          onClick: (e) => {
            e.preventDefault();
            navigate({ to: "/access/request/" + rule.original.id });
          },
        }),
      })}
    </Flex>
  );
};
