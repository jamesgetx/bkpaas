<template>
  <span class="icon-status-container">
    <slot name="icon">
      <i
        class="icon-content"
        :class="mergedIconClass"
        :style="{ color: validColor }"
      ></i>
    </slot>
    <slot name="label">
      <span
        class="label-text"
        :style="{ fontSize: `${fontSize}px` }"
      >
        {{ label || '--' }}
      </span>
    </slot>
  </span>
</template>

<script>
const COLOR_REGEX = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$|^rgb$\d+,\s*\d+,\s*\d+$$/;

export default {
  name: 'IconStatus',
  props: {
    iconClass: {
      type: String,
      validator: (value) => value.startsWith('paasng-'), // 类名格式验证
    },
    iconColor: {
      type: String,
      default: '#3A84FF',
      validator: (value) => COLOR_REGEX.test(value), // 颜色格式验证
    },
    label: {
      type: String,
      default: '',
    },
    fontSize: {
      type: Number,
      default: 12,
    },
  },
  computed: {
    // 合并基础类名和传入的类名
    mergedIconClass() {
      return ['paasng-icon', this.iconClass];
    },
    validColor() {
      return COLOR_REGEX.test(this.iconColor) ? this.iconColor : '#3A84FF';
    },
  },
};
</script>

<style lang="scss" scoped>
.icon-status-container {
  display: inline-flex;
  align-items: center;
  gap: 5px;

  .icon-content {
    --icon-size: 14px;
    --icon-color: v-bind(validColor);

    font-size: var(--icon-size);
    color: var(--icon-color);
    transform: translateY(0px);

    &:hover {
      opacity: 0.8;
    }
  }

  .label-text {
    font-size: 12px;
    color: #313238;
  }
}
</style>
