logical_product_service = "sbntnsgassotn"
network_map = {
  "spoke1" = {
    address_space = ["192.0.0.0/16"]
    subnets = {
      somesbt = {
        prefix = "192.0.0.0/24"
      }
    }
    bgp_community        = null
    ddos_protection_plan = null
    dns_servers          = []
    tags                 = {}
  }
}
