"""
Generates presentations/img/onprem-architecture.png

Active/active two-DC on-prem deployment of Anypoint Flex Gateway with:
  - F5 GSLB cross-DC failover
  - F5 LTM per DC
  - GitOps + Ansible AWX for config push (Local mode) OR outbound mTLS to
    Anypoint Control Plane (Connected mode) - both shown
  - HashiCorp Vault with cross-DC replication
  - Splunk + Datadog observability dual-shipping
  - On-prem MS Integration Stack as the backend

Run:
    python presentations/diagrams/onprem_architecture.py
"""

from diagrams import Diagram, Cluster, Edge
from diagrams.onprem.client import Users
from diagrams.onprem.compute import Server
from diagrams.onprem.monitoring import Splunk, Datadog
from diagrams.onprem.iac import Ansible
from diagrams.onprem.vcs import Github
from diagrams.onprem.security import Vault
from diagrams.generic.network import Router, Firewall
from diagrams.k8s.compute import Pod
from diagrams.saas.identity import Okta
from diagrams.aws.network import APIGateway


graph_attr = {
    "fontsize": "22",
    "bgcolor": "white",
    "pad": "1.0",
    "ranksep": "1.3",
    "nodesep": "1.0",
    "splines": "ortho",
    "compound": "true",
}
node_attr = {"fontsize": "13"}
edge_attr = {"fontsize": "11"}


with Diagram(
    "MuleSoft Flex Gateway - On-Prem Active/Active Across DCs",
    show=False,
    direction="LR",
    filename="presentations/img/onprem-architecture",
    outformat="png",
    graph_attr=graph_attr,
    node_attr=node_attr,
    edge_attr=edge_attr,
):
    ext_users = Users("External users\npartners / web / mobile")
    int_users = Users("Internal users /\nservices")

    with Cluster("Edge (multi-site)"):
        gslb = Router("F5 BIG-IP DNS\nGSLB - cross-DC failover\nTTL 30s")

    with Cluster("Anypoint Control Plane (SaaS - Connected mode only)"):
        anypoint = APIGateway("API Manager\nMonitoring\nAudit\n(outbound HTTPS only)")

    idp = Okta("Enterprise IdP\nEntra / Okta / Ping\nOIDC + JWKS")

    with Cluster("Data Center 1 (active)"):
        fw1 = Firewall("Perimeter FW")
        lb1 = Router("F5 LTM\nVIP 10.1.1.10")
        with Cluster("Flex Gateway pool"):
            fg1a = Pod("flex-dc1-1")
            fg1b = Pod("flex-dc1-2")
        ms1 = Server("MS Integration\nStack\n(BizTalk / Logic Apps)")

    with Cluster("Data Center 2 (active)"):
        fw2 = Firewall("Perimeter FW")
        lb2 = Router("F5 LTM\nVIP 10.2.1.10")
        with Cluster("Flex Gateway pool"):
            fg2a = Pod("flex-dc2-1")
            fg2b = Pod("flex-dc2-2")
        ms2 = Server("MS Integration\nStack\n(replicated)")

    with Cluster("Config + Secrets"):
        gh = Github("Git repo\n(source of truth)")
        ans = Ansible("Ansible AWX\nper-DC push +\nSIGHUP hot reload")
        vault = Vault("HashiCorp Vault\nperformance replication\ncross-DC")

    with Cluster("Observability"):
        splunk = Splunk("Splunk HEC\nlogs + audit")
        dd = Datadog("Datadog APM\nmetrics + traces")

    # ---- Traffic edges ----
    ext_users >> gslb
    int_users >> gslb
    gslb >> Edge(label="health-checked") >> fw1
    gslb >> Edge(label="health-checked") >> fw2
    fw1 >> lb1
    fw2 >> lb2
    lb1 >> fg1a
    lb1 >> fg1b
    lb2 >> fg2a
    lb2 >> fg2b
    fg1a >> ms1
    fg1b >> ms1
    fg2a >> ms2
    fg2b >> ms2
    ms1 >> Edge(style="dashed", label="data replication") >> ms2

    # ---- Connected-mode (outbound only) ----
    fg1a >> Edge(style="dashed", label="policy push\nmTLS HTTPS\n(Connected mode)") >> anypoint
    fg2a >> Edge(style="dashed") >> anypoint

    # ---- Local-mode config replication ----
    gh >> ans
    ans >> Edge(label="parallel push") >> fg1a
    ans >> Edge(label="parallel push") >> fg2a

    # ---- Secrets ----
    vault >> Edge(style="dotted", label="cert + secret\npull at runtime") >> fg1a
    vault >> Edge(style="dotted") >> fg2a

    # ---- JWKS ----
    fg1a >> Edge(style="dotted", color="darkgreen", label="JWKS") >> idp
    fg2a >> Edge(style="dotted", color="darkgreen") >> idp

    # ---- Observability ----
    fg1a >> Edge(style="dashed", color="gray") >> splunk
    fg2a >> Edge(style="dashed", color="gray") >> splunk
    fg1a >> Edge(style="dashed", color="gray") >> dd
    fg2a >> Edge(style="dashed", color="gray") >> dd
