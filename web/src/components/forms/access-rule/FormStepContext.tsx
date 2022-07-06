import React from "react";
import { createCtx } from "../../../utils/context/createCtx";
import { StepActionKind, useSteps } from "./StepsContext";

export interface FormStepContextProps {
  active: boolean;
  showEdit: boolean;
  showClose: boolean;
  showNext: boolean;
  showSubmit: boolean;
  showPreview: boolean;
  next: () => void;
  close: () => void;
  edit: () => void;
}

const [
  useFormStep,
  FormStepContextProvider,
] = createCtx<FormStepContextProps>();

// FormStepcontext is used to give the FormStep number to the formFormStep automatically
const FormStepProvider: React.FC<{
  step: number;
}> = ({ children, step }) => {
  const { stepDispatch, stepState } = useSteps();
  const active = step === stepState.step && !stepState.isReadOnly;
  const isPreviousStep = step < stepState.step;
  const isFinalStep = step === stepState.numberOfSteps - 1;
  return (
    <FormStepContextProvider
      value={{
        active,
        showPreview: !active,
        showEdit: (isPreviousStep || stepState.editMode) && !active,
        showClose: active && stepState.editMode,
        showNext: !stepState.editMode && active && !isFinalStep,
        showSubmit: !stepState.editMode && isFinalStep,
        close: () => stepDispatch({ type: StepActionKind.CLOSE }),
        next: () => stepDispatch({ type: StepActionKind.NEXT }),
        edit: () => stepDispatch({ type: StepActionKind.EDIT, payload: step }),
      }}
    >
      {children}
    </FormStepContextProvider>
  );
};

export { FormStepProvider, useFormStep };
