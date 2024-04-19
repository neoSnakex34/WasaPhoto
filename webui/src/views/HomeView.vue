<script>
export default {
	data: function() {
		return {
			username: localStorage.getItem('username'),
			userId: localStorage.getItem('userId'),
			errormsg: null,
			loading: false,
			some_data: null,
		}
	},
	methods: {
		async refresh() {

			this.loading = true;
			this.errormsg = null;
			try {
				let response = await this.$axios.get("/");
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
		},
		
		doLogout() {
			localStorage.removeItem('userId');
			localStorage.removeItem('username');
			this.$router.push('/login');
		},
		

	},
	mounted() {
		this.refresh()
	}
}
</script>

<template>
	<div>
		<div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h1">{{ this.username }}'s WasaPHOTO</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-4">
					<button type="button" class="btn btn-primary btn-lg" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-danger btn-lg" @click="doLogout">
						Logout
					</button>
				</div>
				<!-- <div class="btn-group me-2">
					<button type="button" class="btn btn-sm" @click="newItem">
						Upload
					</button>
				</div> -->
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>

<style>
</style>
