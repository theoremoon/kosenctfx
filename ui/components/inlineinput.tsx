import React from "react";

type InlineInputProps = Omit<React.ComponentPropsWithRef<"input">, "className">;

const InlineInput = React.forwardRef((props: InlineInputProps, ref) => {
  return (
    <input
      {...props}
      ref={ref}
      className="inline-block w-1/2 bg-transparent color-black border-b focus:border-b-2 mb-px focus:mb-0 border-pink-600 outline-none text-center text-xl"
    />
  );
});

export default InlineInput;
