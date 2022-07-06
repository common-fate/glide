import { ChevronDownIcon } from "@chakra-ui/icons";
import {
  Button,
  ButtonProps,
  chakra,
  Menu,
  MenuButton,
  MenuItemOption,
  MenuList,
  MenuOptionGroup,
  Spinner,
} from "@chakra-ui/react";
import React from "react";
import { useGetUsers } from "../../utils/backend-client/admin/admin";

type Props = {
  setApprovers: React.Dispatch<React.SetStateAction<string[]>>;
} & ButtonProps;

const ApproversInput = ({ setApprovers, ...props }: Props) => {
  const { data, isValidating } = useGetUsers({
    // request: {
    //   baseURL: "http://127.0.0.1:3100",
    //   headers: {
    //     "Content-Type": "application/json",
    //     "Prefer": "code=200, example=ChrisAndJosh",
    //   },
    // },
  });

  return (
    <Menu closeOnSelect={false}>
      <MenuButton
        as={Button}
        rightIcon={isValidating ? <Spinner size="sm" /> : <ChevronDownIcon />}
        variant="outline"
        {...props}
      >
        Select Approvers
      </MenuButton>
      <MenuList minWidth="240px">
        <MenuOptionGroup
          type="checkbox"
          onChange={(approvers) =>
            setApprovers(typeof approvers == "string" ? [approvers] : approvers)
          }
        >
          {data?.users?.map((user) => (
            <MenuItemOption key={user.id} value={user.id} maxW="30ch">
              {user.firstName} {user.lastName}
              <chakra.span
                // w="100%"
                mt={-1}
                color="neutrals.600"
                textAlign="right"
                overflow="hidden"
                textOverflow="ellipsis"
                display="-webkit-box"
                fontSize="sm"
              >
                {user.email}
              </chakra.span>
            </MenuItemOption>
          ))}
        </MenuOptionGroup>
      </MenuList>
    </Menu>
  );
};

export default ApproversInput;
