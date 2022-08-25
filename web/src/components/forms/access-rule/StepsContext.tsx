import React, { Children, useReducer } from "react";
import { createCtx } from "../../../utils/context/createCtx";
import { FormStepProvider } from "./FormStepContext";

// An enum with all the types of actions to use in our reducer
export enum StepActionKind {
  EDIT = "EDIT",
  NEXT = "NEXT",
  CLOSE = "CLOSE",
}

// An interface for our actions
export interface StepAction {
  type: StepActionKind;
  payload?: number;
}

// An interface for our state
export interface StepState {
  editMode: boolean;
  isReadOnly: boolean;
  step: number;
  numberOfSteps: number;
}

// Our reducer function that uses a switch statement to handle our actions
function stepReducer(state: StepState, action: StepAction): StepState {
  const { type, payload } = action;
  switch (type) {
    case StepActionKind.EDIT:
      if (payload !== undefined) {
        return {
          ...state,
          step: payload,
        };
      }
      return state;

    case StepActionKind.NEXT:
      if (state.step !== -1 && state.step < state.numberOfSteps - 1) {
        return {
          ...state,
          step: state.step + 1,
        };
      }
      break;
    case StepActionKind.CLOSE:
      return {
        ...state,
        step: -1,
      };
    default:
      return state;
  }
  return state;
}

export interface StepsContextProps {
  stepState: StepState;
  stepDispatch: React.Dispatch<StepAction>;
}

const [useSteps, StepsContextProvider] = createCtx<StepsContextProps>();

const StepsProvider: React.FC<{
  isEditMode?: boolean;
  isReadOnly?: boolean;
  children?: React.ReactElement[];
}> = ({ isEditMode, children, isReadOnly }) => {
  const editMode = !!isEditMode;
  const readOnlyMode = !!isReadOnly;
  const [state, dispatch] = useReducer(stepReducer, {
    // -1 means that no step is active
    step: editMode ? -1 : 0,
    numberOfSteps: children ? children.length : 0,
    editMode: editMode,
    isReadOnly: readOnlyMode,
  });

  return (
    <StepsContextProvider
      value={{
        stepState: state,
        stepDispatch: dispatch,
      }}
    >
      {/* Each form step is wrapped in a step context which tells it which step index it is */}
      {Children.map(children, (child, index) => {
        return (
          child && <FormStepProvider step={index}>{child}</FormStepProvider>
        );
      })}
    </StepsContextProvider>
  );
};

export { useSteps, StepsProvider };
