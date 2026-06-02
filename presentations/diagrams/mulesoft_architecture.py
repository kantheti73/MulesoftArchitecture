"""
Generates presentations/img/mulesoft-architecture.png

Run:
    # one-time
    winget install Graphviz.Graphviz
    pip install --user diagrams
    # whenever the diagram changes
    python presentations/diagrams/mulesoft_architecture.py
"""

from diagrams import Diagram, Cluster, Edge
from diagrams.aws.network import APIGateway, DirectConnect, TransitGateway, SiteToSiteVpn
from diagrams.azure.network import (
    ExpressrouteCircuits, VirtualNetworks, LoadBalancers, ApplicationGateway,
    VirtualNetworkGateways,
)
from diagrams.azure.identity import ActiveDirectory
from diagrams.azure.security import KeyVaults
from diagrams.azure.compute import ContainerInstances
from diagrams.azure.web import APIConnections
from diagrams.azure.general import Resource, Subscriptions
from diagrams.onprem.client import Users
from diagrams.onprem.compute import Server
from diagrams.onprem.monitoring import Splunk, Datadog
from diagrams.generic.network import Router, Firewall


graph_attr = {
    "fontsize": "22",
    "bgcolor": "white",
    "pad": "1.0",
    "ranksep": "1.2",
    "nodesep": "0.9",
    "splines": "ortho",
    "compound": "true",
}
node_attr = {"fontsize": "15"}
edge_attr = {"fontsize": "12"}


with Diagram(
    "MuleSoft Flex Gateway - Federated, Private, Dual-Listener",
    show=False,
    direction="LR",
    filename="presentations/img/mulesoft-architecture",
    outformat="png",
    graph_attr=graph_attr,
    node_attr=node_attr,
    edge_attr=edge_attr,
):
    ext_users = Users("External users\npartners / web / mobile")
    int_users = Users("Internal users /\nservices")
    idp = ActiveDirectory("SSP IdP\n(Entra / Okta)\nOAuth + OIDC")

    with Cluster("Anypoint Control Plane (SaaS)"):
        anypoint = APIGateway("API Manager\nExchange\nMonitoring\nAudit Log")

    with Cluster("MuleSoft Private Space  (Azure centralus)"):
        with Cluster("External Ingress"):
            waf = ApplicationGateway("Anypoint\nSecurity Edge\nWAF + DDoS")
            dlb = LoadBalancers("Anypoint DLB\napi.yourco.com\nTLS termination")
        with Cluster("Internal Ingress"):
            intlb = LoadBalancers("Internal LB\napi-internal.yourco.local")
        with Cluster("Flex Gateway (n=2, multi-AZ)"):
            fg1 = ContainerInstances("Flex Gateway\nreplica 1")
            fg2 = ContainerInstances("Flex Gateway\nreplica 2")
        psvnet = VirtualNetworks("Private Space VNet\n10.50.0.0/22")

    ergw = VirtualNetworkGateways("ExpressRoute\nGateway")
    er = ExpressrouteCircuits("ExpressRoute\nPrivate Peering\n10 Gbps - primary")
    vpn = SiteToSiteVpn("IPsec VPN\nbackup over ER\nMicrosoft Peering")

    with Cluster("On-Premises Data Center"):
        fw = Firewall("Perimeter FW")
        f5 = Router("F5 BIG-IP\nLTM")
        ms = Server("Microsoft Integration\nStack\n(BizTalk / Logic Apps /\nEDI engine)")
        backend = Server("Downstream apps\n+ DBs")

    with Cluster("Observability"):
        splunk = Splunk("Splunk HEC\nlogs + audit")
        dd = Datadog("Datadog\nmetrics + traces")

    # External path
    ext_users >> waf >> dlb >> fg1
    dlb >> fg2

    # Internal path
    int_users >> intlb >> fg1
    intlb >> fg2

    # SaaS control plane (outbound only, mTLS)
    fg1 >> Edge(style="dashed", label="policy + telemetry\nmTLS outbound") >> anypoint
    fg2 >> Edge(style="dashed") >> anypoint

    # JWT validation (gateway -> IdP)
    fg1 >> Edge(style="dotted", color="darkgreen", label="JWKS / OIDC") >> idp
    fg2 >> Edge(style="dotted", color="darkgreen") >> idp

    # Private connectivity
    fg1 >> psvnet
    fg2 >> psvnet
    psvnet >> Edge(label="VNet peering") >> ergw
    ergw >> Edge(label="BGP / Private Peering") >> er
    ergw >> Edge(style="dashed", label="failover") >> vpn
    er >> fw
    vpn >> fw
    fw >> f5 >> ms >> backend

    # Observability dual-shipping
    fg1 >> Edge(style="dashed", color="gray") >> splunk
    fg1 >> Edge(style="dashed", color="gray") >> dd
