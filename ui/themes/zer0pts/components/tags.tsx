import { HStack, Tag } from "@chakra-ui/react";

interface TagsProps {
  tags: string[];
}

const Tags = ({ tags }: TagsProps) => {
  return (
    <HStack>
      {tags.map((t) => (
        <Tag
          key={t}
          colorScheme="blackAlpha"
          variant="solid"
          sx={{ color: "black", backgroundColor: "#0000001a" }}
        >
          {t}
        </Tag>
      ))}
    </HStack>
  );
};

export default Tags;
