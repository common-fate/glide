import { ChevronDownIcon } from "@chakra-ui/icons";
import {
  Menu,
  MenuButton,
  Button,
  MenuList,
  MenuOptionGroup,
  MenuItemOption,
  Spinner,
  ButtonProps,
  MenuOptionGroupProps,
} from "@chakra-ui/react";
import React from "react";
import { useGetGroups } from "../../utils/backend-client/admin/admin";

type Props = {
  onChange: (...event: any[]) => void;
  // setGroupIds: React.Dispatch<React.SetStateAction<string[]>>;
} & MenuOptionGroupProps;

/**
 * To be initialised with a state setter like so:
 * const [groupIds, setGroupIds] = React.useState<string[]>([]);
 */

const GroupInput = ({ ...props }: Props) => {
  const { data, isValidating } = useGetGroups({});
  return (
    <Menu closeOnSelect={false}>
      <MenuButton
        as={Button}
        rightIcon={isValidating ? <Spinner size="sm" /> : <ChevronDownIcon />}
        variant="outline"
        bg="white"
        maxW={{ md: "3xl" }}
        placeholder="Placeholder"
      >
        Select groups
      </MenuButton>
      <MenuList minWidth="240px">
        {!isValidating && data?.groups && (
          <MenuOptionGroup type="checkbox">
            {data?.groups?.map((group) => (
              <MenuItemOption key={group.name} value={group.name}>
                {group.name}
              </MenuItemOption>
            ))}
          </MenuOptionGroup>
        )}
      </MenuList>
    </Menu>
  );
};

export default GroupInput;
