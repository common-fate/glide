import { ChevronDownIcon } from "@chakra-ui/icons";
import {
  Menu,
  MenuButton,
  Button,
  MenuList,
  MenuOptionGroup,
  MenuItemOption,
} from "@chakra-ui/react";
import React from "react";
import { AdminListGroupsSource } from "../../utils/backend-client/types";

export const GroupsFilterMenu: React.FC<{
  source: AdminListGroupsSource | undefined;
  onChange: (source: AdminListGroupsSource | undefined) => void;
}> = ({ source, onChange }) => {
  return (
    <Menu>
      <MenuButton
        as={Button}
        rightIcon={<ChevronDownIcon />}
        variant="ghost"
        size="sm"
      >
        {source === "INTERNAL"
          ? "Internal Only"
          : source === "EXTERNAL"
          ? "External Only"
          : "All"}
      </MenuButton>
      <MenuList>
        <MenuOptionGroup
          defaultValue="all"
          title="View option"
          type="radio"
          onChange={(e) => {
            switch (e) {
              case "int":
                onChange(AdminListGroupsSource.INTERNAL);
                break;
              case "ext":
                onChange(AdminListGroupsSource.EXTERNAL);
                break;

              default:
                onChange(undefined);
            }
          }}
        >
          <MenuItemOption value="all">All</MenuItemOption>
          <MenuItemOption value="int">Internal Only</MenuItemOption>
          <MenuItemOption value="ext">External Only</MenuItemOption>
        </MenuOptionGroup>
      </MenuList>
    </Menu>
  );
};
