import type { Action } from "svelte/action";

export type SyncDialogParams = {
  open: boolean;
  onPresent?: () => void;
};

/** Keeps a native <dialog> in sync with a boolean without $effect (avoids effect_update_depth_exceeded). */
export const syncDialog: Action<HTMLDialogElement, SyncDialogParams> = (node, params) => {
  let wasOpen = false;

  const apply = ({ open, onPresent }: SyncDialogParams) => {
    if (open && !wasOpen) {
      onPresent?.();
    }
    wasOpen = open;
    if (open) {
      if (!node.open) {
        node.showModal();
      }
    } else if (node.open) {
      node.close();
    }
  };

  apply(params);

  return {
    update(next) {
      apply(next);
    },
  };
};
