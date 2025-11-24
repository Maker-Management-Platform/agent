import React from 'react';
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
} from '@mantine/core';
import { useForm } from '@mantine/form';
import { useMutation, useQueryClient } from '@tanstack/react-query';

import { type Asset, type AssetPatch } from '../../models';
import { patchAsset } from '../../fetchers/getAsset';

/**
 *
 * @param param0
 * @param param0.asset
 * @param param0.onClose
 * @returns
 */
function AssetEditForm({
  asset,
  onClose,
}: {
  readonly asset: Asset;
  readonly onClose: () => void;
}) {
  const queryClient = useQueryClient();
  const mutation = useMutation({
    mutationFn: patchAsset,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['asset', asset.ID] }).catch(console.error);
      onClose();
    },
  });
  const form = useForm<AssetPatch>({
    initialValues: {
      Label: asset.Label,
      Description: asset.Description ?? '',
      Tags: asset.Tags ?? [],
    } as AssetPatch,
    validate: {
      Label: (value) => ((!value || value.length > 0) ? undefined : 'Label is required'),
      Description: (value) => ((!value || value.length > 0) ? undefined : 'Description is required'),
    },
  });
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const formSubmit = (values: AssetPatch): Promise<any> => (mutation.mutateAsync(values));
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
            key={form.key('Label')}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...form.getInputProps('Label')}
          />
          <Textarea
            mt="md"
            label="Description"
            placeholder="A cool thing to print"
            key={form.key('Description')}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...form.getInputProps('Description')}
          />
        </Fieldset>
        <Fieldset legend="Asset Context">
          <TagsInput
            label="Tags"
            splitChars={[',', ' ', '|']}
            placeholder="tag1, tag2"
            key={form.key('Tags')}
            // eslint-disable-next-line react/jsx-props-no-spreading
            {...form.getInputProps('Tags')}
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

export default AssetEditForm;
