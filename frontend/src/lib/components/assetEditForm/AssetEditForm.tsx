import {
  Button,
  CloseButton,
  Divider,
  Fieldset,
  Group,
  SimpleGrid,
  TagsInput,
  Textarea,
  TextInput,
} from "@mantine/core";
import { Asset, AssetPatch } from "../../models";
import { useForm } from "@mantine/form";
import { patchAsset } from "../../fetchers/getAsset";
import { useMutation, useQueryClient } from "@tanstack/react-query";

export function AssetEditForm({
  asset,
  onClose,
}: {
  asset: Asset;
  onClose: () => void;
}) {
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: patchAsset,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["asset", asset.ID] });
      onClose();
    },
  });
  const form = useForm<AssetPatch>({
    initialValues: {
      Label: asset.Label,
      Description: asset.Description ? asset.Description : "",
      Tags: asset.Tags ? asset.Tags : [],
    } as AssetPatch,
    validate: {
      Label: (value) =>
        !value || value.length > 0 ? undefined : "Label is required",
      Description: (value) =>
        !value || value.length > 0 ? undefined : "Description is required",
    },
  });
  const formSubmit = (values: AssetPatch): Promise<any> => {
    values.ID = asset.ID;
    return mutation.mutateAsync(values);
  };
  return (
    <form onSubmit={form.onSubmit(formSubmit)} onReset={form.reset}>
      <Group grow>
        <Divider label="Edit asset" />
        <CloseButton onClick={() => onClose()} />
      </Group>
      <SimpleGrid cols={3}>
        <Fieldset legend="Asset details">
          <TextInput
            label="Name"
            placeholder="An asset"
            key={form.key("Label")}
            {...form.getInputProps("Label")}
          />
          <Textarea
            mt="md"
            label="Description"
            placeholder="A cool thing to print"
            key={form.key("Description")}
            {...form.getInputProps("Description")}
          />
        </Fieldset>
        <Fieldset legend="Asset Context">
          <TagsInput
            label="Tags"
            splitChars={[",", " ", "|"]}
            placeholder="tag1, tag2"
            key={form.key("Tags")}
            {...form.getInputProps("Tags")}
          />
        </Fieldset>
      </SimpleGrid>
      <Group align="flex-end" mt="lg">
        <Button type="submit">Save</Button>
        <Button
          variant="light"
          onClick={() => {
            form.reset();
            onClose();
          }}
        >
          Cancel
        </Button>
      </Group>
    </form>
  );
}

