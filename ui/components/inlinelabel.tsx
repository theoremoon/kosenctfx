import React from "react";

type InlineLabelProps = Omit<
  Rect.ComponentPropsWithoutRef<"label">,
  "className"
>;

const InlineLabel = (props: InlineLabelProps) => {
  return <label {...props} className="inline-block w-1/2 font-bold text-lg" />;
};

export default InlineLabel;
