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
        {statusMenuOption(status)}
      </MenuButton>
      <MenuList>
        <MenuOptionGroup
          defaultValue="all"
          title="View option"
          type="radio"
          value={status}
          onChange={(e) => {
            switch (typeof e) {
              case "string":
                onChange(e as RequestStatus);
                break;
              default:
                onChange(undefined);
            }
          }}
        >
          <MenuItemOption value={"ALL"}>All</MenuItemOption>
          <MenuItemOption value={RequestStatus.PENDING}>
            {statusMenuOption(RequestStatus.PENDING)}
          </MenuItemOption>
          <MenuItemOption value={RequestStatus.ACTIVE}>
            {statusMenuOption(RequestStatus.ACTIVE)}
          </MenuItemOption>
          <MenuItemOption value={RequestStatus.COMPLETE}>
            {statusMenuOption(RequestStatus.COMPLETE)}
          </MenuItemOption>
          <MenuItemOption value={RequestStatus.CANCELLED}>
            {statusMenuOption(RequestStatus.CANCELLED)}
          </MenuItemOption>
          <MenuItemOption value={RequestStatus.REVOKED}>
            {statusMenuOption(RequestStatus.REVOKED)}
          </MenuItemOption>
        </MenuOptionGroup>
      </MenuList>
    </Menu>
  );
};

const statusMenuOption = (
  status: RequestStatus | "ALL" | undefined
): string => {
  return status === "PENDING"
    ? "Pending only"
    : status === "ACTIVE"
    ? "Active only"
    : status === "COMPLETE"
    ? "Complete only"
    : status === "CANCELLED"
    ? "Cancelled only"
    : status === "REVOKED"
    ? "Revoked only"
    : "All";
};
