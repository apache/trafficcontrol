import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import { CDNService, ProfileService } from 'src/app/api';
import { NavigationService } from 'src/app/shared/navigation/navigation.service';
import { ResponseCDN, ResponseProfile } from 'trafficops-types';

@Component({
  selector: 'tp-profile-detail',
  templateUrl: './profile-detail.component.html',
  styleUrls: ['./profile-detail.component.scss']
})
export class ProfileDetailComponent implements OnInit {
  public new = false;
  public profile!: ResponseProfile;
  public cdns!: ResponseCDN[];
  public form = new FormGroup({
    cdn: new FormControl(0, { nonNullable: true }),
    description: new FormControl("", { nonNullable: true }),
    name: new FormControl("", { nonNullable: true }),
    routingDisabled: new FormControl(false, { nonNullable: true }),
    type: new FormControl("", { nonNullable: true })
  });
  types = [
    { value: 'ATS_PROFILE' },
    { value: 'TR_PROFILE' },
    { value: 'TM_PROFILE' },
    { value: 'TS_PROFILE' },
    { value: 'TP_PROFILE' },
    { value: 'INFLUXDB_PROFILE' },
    { value: 'RIAK_PROFILE' },
    { value: 'SPLUNK_PROFILE' },
    { value: 'DS_PROFILE' },
    { value: 'ORG_PROFILE' },
    { value: 'KAFKA_PROFILE' },
    { value: 'LOGSTASH_PROFILE' },
    { value: 'ES_PROFILE' },
    { value: 'UNK_PROFILE' },
    { value: 'GROVE_PROFILE' }
  ];

  constructor(
    private readonly route: ActivatedRoute,
    private readonly navSvc: NavigationService,
    private api: ProfileService,
    private cdnService: CDNService
  ) { }

  public async ngOnInit(): Promise<void> {
    this.cdns = await this.cdnService.getCDNs();
    const id = this.route.snapshot.paramMap.get("id");

    if (id === null) {
      throw new Error("missing required route parameter 'id'");
    } else if (id === "new") {
      this.new = true;
      this.navSvc.headerTitle.next("New Profile");
    } else {
      const numID = parseInt(id, 10);
      if (Number.isNaN(numID)) {
        throw new Error(`route parameter 'id' was non-number:  ${{ id }}`);
      } else {
        this.profile = await this.api.getProfiles(Number(id));
        this.navSvc.headerTitle.next(`Profile: ${this.profile.name}`);
        this.form.patchValue(this.profile);
      }
    }
  }

  onSubmit() {
    console.log(this.form.value);
  }
}
