import { Button } from "@chakra-ui/react";
import { useState } from "react";
import { adminPostApiV1IdentitySync } from "../utils/backend-client/default/default";

interface Props {
  onSync?: () => void;
}
export const SyncUsersAndGroupsButton: React.FC<Props> = ({ onSync }) => {
  const [isSyncing, setIsSyncing] = useState(false);
  const sync = async () => {
    try {
      setIsSyncing(true);
      await adminPostApiV1IdentitySync();
      onSync?.();
    } finally {
      setIsSyncing(false);
    }
  };
  return (
    <Button size="sm" variant="ghost" onClick={sync} isLoading={isSyncing}>
      Sync Users and Groups
    </Button>
  );
};
