import React, { Component, ReactNode } from "react";
import UnhandledError from "./errPage";

interface Props {
  children?: ReactNode;
}

interface State {
  error: Error | undefined;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    error: undefined,
  };

  public static getDerivedStateFromError(error: Error): State {
    // Update state so the next render will show the fallback UI.
    return { error: error };
  }

  public render() {
    if (this.state.error != undefined) {
      return <UnhandledError error={this.state.error} />;
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
