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
import { RequestStatus } from "../../utils/backend-client/types";

export const RequestsFilterMenu: React.FC<{
  status: RequestStatus | undefined;
  onChange: (status: RequestStatus | undefined) => void;
}> = ({ status, onChange }) => {
  return (
    <Menu>
      <MenuButton
        as={Button}
        rightIcon={<ChevronDownIcon />}
        variant="ghost"
        size="sm"
      >
        {status === "PENDING"
          ? "Pending only"
          : status === "DECLINED"
          ? "Declined only"
          : status === "APPROVED"
          ? "Approved only"
          : status === "CANCELLED"
          ? "Cancelled only"
          : "All"}
      </MenuButton>
      <MenuList>
        <MenuOptionGroup
          defaultValue="all"
          title="View option"
          type="radio"
          onChange={(e) => {
            switch (e) {
              case "pend":
                onChange(RequestStatus.PENDING);
                break;
              case "den":
                onChange(RequestStatus.DECLINED);
                break;
              case "apr":
                onChange(RequestStatus.APPROVED);
                break;
              case "can":
                onChange(RequestStatus.CANCELLED);
                break;
              default:
                onChange(undefined);
            }
          }}
        >
          <MenuItemOption value="all">All</MenuItemOption>
          <MenuItemOption value="pend">Pending only</MenuItemOption>
          <MenuItemOption value="den">Declined only</MenuItemOption>
          <MenuItemOption value="apr">Approved only</MenuItemOption>
          <MenuItemOption value="can">Cancelled only</MenuItemOption>
        </MenuOptionGroup>
      </MenuList>
    </Menu>
  );
};
