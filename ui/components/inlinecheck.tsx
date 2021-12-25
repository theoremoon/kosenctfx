import React from "react";

type InlineCheckProps = Omit<
  React.ComponentPropsWithRef<"input">,
  "className",
  "type"
>;

const InlineCheck = React.forwardRef((props: InlineCheckProps, ref) => {
  return <input ref={ref} type="checkbox" className="w-6 h-6" />;
});

export default InlineCheck;
