import React from "react";

import { Link, Text } from "@chakra-ui/react";
import { useAdminListProviders } from "../utils/backend-client/admin/admin";
import { OnboardingCard } from "./OnboardingCard";
import { WarningIcon } from "./icons/Icons";
const ProviderSetupNotice: React.FC = () => {
  const { data, error } = useAdminListProviders();

  if (error) {
    return (
      <OnboardingCard
        leftIcon={<WarningIcon />}
        title="Oops, there was an error fetching your providers, check that your deployment is correctly configured"
      ></OnboardingCard>
    );
  }
  if (data?.length == 0) {
    return (
      <OnboardingCard
        leftIcon={<WarningIcon />}
        title="It looks like you don't have any providers configured yet"
      >
        <Text>
          Before you can create an access rule, you need to setup your first
          access provider. Use the{" "}
          <Link href="/admin/providers/setup">interactive setup guides</Link> to
          get started!
        </Text>
      </OnboardingCard>
    );
  }
  return null;
};

export default ProviderSetupNotice;
