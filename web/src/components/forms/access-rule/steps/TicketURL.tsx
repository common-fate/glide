import {
    Flex,
    FormControl,
    FormErrorMessage,
    FormHelperText,
    FormLabel,
    HStack,
    Switch,
    Text,
    VStack,
    Wrap,
    WrapItem,
} from "@chakra-ui/react";
import React from "react";
import { useFormContext } from "react-hook-form";
import { useAdminGetGroup } from "../../../../utils/backend-client/admin/admin";
import { GroupSelect, UserSelect } from "../components/Select";
import { AccessRuleFormData } from "../CreateForm";

import { FormStep } from "./FormStep";
import { CFAvatar } from "../../../CFAvatar";
import {GroupDisplay} from "./Approval";

export const TicketURLStep: React.FC = () => {
    console.log("Hello from here!!!")
    const methods = useFormContext();
    const ticketURL = methods.watch("ticketURL");
    // If approval is required, then at least one user or one group needs to be set
    // const approverRequired =
    //     !!approval?.required &&
    //     !(approval?.groups?.length > 0 || approval?.users?.length > 0);
    return (
        <FormStep
            heading="Ticket URL"
            subHeading="Ticket URL"
            fields={[]}
            hideNext={true}
            preview={<TicketURLPreview />}
        >
            <>
                <FormControl>
                    <FormLabel htmlFor="ticketURL">
                        <HStack>
                            <Switch
                                id="requires-ticketURL-button"
                                bg="neutrals.0"
                                {...methods.register("ticketURL", {})}
                            />
                            <Text textStyle={"Body/Medium"}>Ticket URL required</Text>
                        </HStack>
                    </FormLabel>
                </FormControl>
            </>
        </FormStep>
    );
};

const TicketURLPreview: React.FC = () => {
    const methods = useFormContext();
    const required = methods.watch("ticketURL")
    if (required)
        return <Text w="100%">Ticket URL required</Text>;
    else
        return <Text w="100%">No ticket URL required</Text>;
};