<script>
export default {
	props: ['msg'],
	methods: {
		closeErrorMsg() {
			this.$emit('update:msg', null);
		},
	}
}

</script>

<template>
	<div class="alert alert-danger" role="alert">
		{{ msg }}
	</div>
</template>

<style>
</style>
