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
import { GroupSource } from "../../utils/backend-client/types";

export const GroupsFilterMenu: React.FC<{
  source: GroupSource | undefined;
  onChange: (source: GroupSource | undefined) => void;
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
          : // : source === "AZURE"
            // ? "Azure only"
            // : source === "GOOGLE"
            // ? "Google only"
            // : source === "AWS"
            // ? "AWS only"
            "All"}
      </MenuButton>
      <MenuList>
        <MenuOptionGroup
          defaultValue="all"
          title="View option"
          type="radio"
          onChange={(e) => {
            switch (e) {
              case "int":
                onChange(GroupSource.INTERNAL);
                break;
              // case "az":
              //   onChange(GroupSource.AZURE);
              //   break;
              // case "go":
              //   onChange(GroupSource.GOOGLE);
              //   break;
              // case "aws":
              //   onChange(GroupSource.AWS);
              //   break;
              default:
                onChange(undefined);
            }
          }}
        >
          <MenuItemOption value="all">All</MenuItemOption>
          <MenuItemOption value="int">Internal Only</MenuItemOption>
          {/* <MenuItemOption value="az">Azure only</MenuItemOption>
          <MenuItemOption value="go">Google only</MenuItemOption>
          <MenuItemOption value="aws">AWS only</MenuItemOption> */}
        </MenuOptionGroup>
      </MenuList>
    </Menu>
  );
};
