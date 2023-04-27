import { ArrowBackIcon, CheckCircleIcon, SettingsIcon } from "@chakra-ui/icons";
import {
  Box,
  Button,
  Center,
  Container,
  Flex,
  HStack,
  Input,
  Spinner,
  Stack,
  TabPanel,
  TabPanels,
  Tabs,
  Text,
  Textarea,
  Tooltip,
  useBoolean,
  useEventListener,
  VStack,
} from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-location";
import Counter from "../components/Counter";
import FieldsCodeBlock from "../components/FieldsCodeBlock";
import { ProviderIcon, ShortTypes } from "../components/icons/providerIcon";
import { UserLayout } from "../components/Layout";
import {
  userListEntitlementTargets,
  userPostRequests,
  userRequestPreflight,
  useUserListEntitlements,
} from "../utils/backend-client/default/default";
import {
  Preflight,
  Target,
  UserListEntitlementTargetsParams,
} from "../utils/backend-client/types";
import { Command as CommandNew } from "../utils/cmdk";

const Search2: React.FC = () => {
  return (
    <UserLayout>
      <Container>
        <VStack mt={10}>
          <Input />
        </VStack>
      </Container>
    </UserLayout>
  );
};

export default Search2;
