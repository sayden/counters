import { Icon } from '@astrojs/starlight/components';

export default function Button({ href, children }) {
  return (
    <button className="link-button p-[1ch]" type="button">
      <div class="flex gap-[1ch] items-center">
        {children}
      </div>
    </button>
  )
}
