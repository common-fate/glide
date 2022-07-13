import React from "react";

import { Link, Text } from "@chakra-ui/react";
import { useListProviders } from "../utils/backend-client/default/default";
import { OnboardingCard } from "./OnboardingCard";
import { WarningIcon } from "./icons/Icons";
const ProviderSetupNotice: React.FC = () => {
  const { data, error } = useListProviders();

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
          Before you can create an access rule, follow our{" "}
          <Link
            isExternal
            href="https://docs.commonfate.io/granted-approvals/getting-started/acess-provider"
          >
            getting started guide
          </Link>{" "}
          to setup your first provider now!
        </Text>
      </OnboardingCard>
    );
  }
  return null;
};

export default ProviderSetupNotice;
