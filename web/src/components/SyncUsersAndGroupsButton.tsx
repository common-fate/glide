import { RepeatIcon } from "@chakra-ui/icons";
import { Button } from "@chakra-ui/react";
import { useState } from "react";
import { adminIdentitySync } from "../utils/backend-client/admin/admin";

interface Props {
  onSync?: () => void;
}
export const SyncUsersAndGroupsButton: React.FC<Props> = ({ onSync }) => {
  const [isSyncing, setIsSyncing] = useState(false);
  const sync = async () => {
    try {
      setIsSyncing(true);
      await adminIdentitySync();
      onSync?.();
    } finally {
      setIsSyncing(false);
    }
  };
  return (
    <Button
      leftIcon={<RepeatIcon />}
      size="sm"
      variant="ghost"
      onClick={sync}
      isLoading={isSyncing}
      ml="60%"
    >
      Sync Users and Groups
    </Button>
  );
};
